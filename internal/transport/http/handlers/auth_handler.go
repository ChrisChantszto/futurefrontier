package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/config"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuth(r fiber.Router, db *mongo.Database, cfg config.Config, log *zap.Logger) {
	svc := service.NewAuthService(db)
	smtpService := service.NewSMTPService(cfg)
	otpService := service.NewOTPService(db, cfg, smtpService)

	group := r.Group("/auth")

	// GET /auth/check-first-user - Check if first user exists
	group.Get("/check-first-user", func(c *fiber.Ctx) error {
		ctx := c.Context()
		count, err := svc.CountUsers(ctx)
		if err != nil {
			log.Error("Failed to count users", zap.Error(err))
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		// Standard envelope + keep top-level field for compatibility
		return JSONSuccessWithExtra(c, "OK", nil, fiber.Map{
			"exists": count > 0,
		})
	})

	group.Post("/init-first-user", func(c *fiber.Ctx) error {
		var body struct{ Email, Password string }
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}

		ctx := c.Context()
		count, err := svc.CountUsers(ctx)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		if count > 0 && body.Email != cfg.SuperAdminEmail {
			return fiber.NewError(fiber.StatusForbidden, "already initialized")
		}
		_, err = svc.CreateUser(ctx, body.Email, body.Password, []models.Role{models.RoleSuperAdmin})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "create user failed")
		}
		return JSONSuccessWithExtra(c, "User initialized", nil, fiber.Map{"ok": true})
	})

	group.Post("/login", func(c *fiber.Ctx) error {
		var body struct{ Email, Password string }
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		u, err := svc.FindUserByEmail(c.Context(), body.Email)
		if err != nil || !svc.VerifyPassword(u, body.Password) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}

		accessToken, _ := createJWT(cfg.JWTSecret, u, 15*time.Minute)
		refreshToken, _ := createJWT(cfg.JWTRefreshSecret, u, 7*24*time.Hour)

		setCookie(c, "access", accessToken, 15*time.Minute, cfg)
		setCookie(c, "refresh", refreshToken, 7*24*time.Hour, cfg)

		return JSONSuccessWithExtra(c, "User logged in successfully", nil, fiber.Map{"ok": true})
	})

	group.Post("/logout", func(c *fiber.Ctx) error {
		clearCookie(c, "access", cfg)
		clearCookie(c, "refresh", cfg)
		return JSONSuccessWithExtra(c, "Logged out", nil, fiber.Map{"ok": true})
	})

	group.Post("/refresh", func(c *fiber.Ctx) error {
		rt := string(c.Cookies("refresh"))
		if rt == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "no refresh")
		}
		claims, err := parseJWT(cfg.JWTRefreshSecret, rt)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh")
		}
		u := &models.User{Email: claims["email"].(string), Roles: []models.Role{}}
		at, _ := createJWT(cfg.JWTSecret, u, 15*time.Minute)
		setCookie(c, "access", at, 15*time.Minute, cfg)
		return JSONSuccessWithExtra(c, "Token refreshed", nil, fiber.Map{"ok": true})
	})

	group.Post("/forgot-password", func(c *fiber.Ctx) error {
		var body struct{ Email string }
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		u, err := svc.FindUserByEmail(c.Context(), body.Email)
		if err != nil {
			// don’t reveal existence
			return JSONSuccessWithExtra(c, "If an account exists, we've sent a reset instruction", nil, fiber.Map{"ok": true})
		}
		token, _ := svc.GenerateResetToken(c.Context(), u.ID, 30*time.Minute)
		// TODO: send email via EmailService; if SMTP not set, return token for dev
		return JSONSuccessWithExtra(c, "Reset token generated", fiber.Map{"devToken": token}, fiber.Map{"ok": true})
	})

	group.Post("/reset-password", func(c *fiber.Ctx) error {
		var body struct {
			Token    string
			Password string
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		rt, err := svc.ConsumeResetToken(c.Context(), body.Token)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid token")
		}
		// update user password
		_, err = db.Collection("users").UpdateByID(c.Context(), rt.UserID, bson.M{"$set": bson.M{
			"passwordHash": hashPassword(body.Password),
		}})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "update failed")
		}
		return JSONSuccessWithExtra(c, "Password reset successfully", nil, fiber.Map{"ok": true})
	})

	// OTP endpoints
	otpGroup := group.Group("/otp")

	// POST /auth/otp/request
	otpGroup.Post("/request", func(c *fiber.Ctx) error {
		var body struct {
			Email   string `json:"email"`
			Purpose string `json:"purpose"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}

		// Validate email
		email := strings.TrimSpace(body.Email)
		if email == "" {
			return fiber.NewError(fiber.StatusBadRequest, "email is required")
		}

		// Validate purpose
		purpose := models.OTPPurpose(body.Purpose)
		if !purpose.IsValid() {
			return fiber.NewError(fiber.StatusBadRequest, "invalid purpose")
		}

		// Request OTP
		requestID, err := otpService.RequestOTP(c.Context(), email, purpose)
		if err != nil {
			log.Error("OTP request failed", zap.Error(err), zap.String("email", email))
			// Return generic message to avoid enumeration
			return JSONSuccessWithExtra(c, "If an account exists, we've sent a code to your email", nil, nil)
		}

		// Success: include requestId in both data and top-level for compatibility
		return JSONSuccessWithExtra(c, "OTP sent", fiber.Map{"requestId": requestID}, fiber.Map{"requestId": requestID})
	})

	// POST /auth/otp/verify
	otpGroup.Post("/verify", func(c *fiber.Ctx) error {
		var body struct {
			Email     string `json:"email"`
			Code      string `json:"code"`
			RequestID string `json:"requestId,omitempty"`
			Purpose   string `json:"purpose"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}

		// Validate inputs
		email := strings.TrimSpace(body.Email)
		code := strings.TrimSpace(body.Code)
		if email == "" || code == "" {
			return fiber.NewError(fiber.StatusBadRequest, "email and code are required")
		}

		// Validate purpose
		purpose := models.OTPPurpose(body.Purpose)
		if !purpose.IsValid() {
			return fiber.NewError(fiber.StatusBadRequest, "invalid purpose")
		}

		// Verify OTP
		otpDoc, err := otpService.VerifyOTP(c.Context(), email, code, body.RequestID, purpose)
		if err != nil {
			log.Error("OTP verification failed", zap.Error(err), zap.String("email", email))
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired code")
		}

		// For login purpose, create JWT tokens
		if purpose == models.OTPPurposeLogin {
			// Find user by email
			user, err := svc.FindUserByEmail(c.Context(), email)
			if err != nil {
				log.Error("User not found after OTP verification", zap.Error(err), zap.String("email", email))
				return fiber.NewError(fiber.StatusUnauthorized, "user not found")
			}

			// Create tokens
			accessToken, _ := createJWT(cfg.JWTSecret, user, 15*time.Minute)
			refreshToken, _ := createJWT(cfg.JWTRefreshSecret, user, 7*24*time.Hour)

			// Set cookies
			setCookie(c, "access", accessToken, 15*time.Minute, cfg)
			setCookie(c, "refresh", refreshToken, 7*24*time.Hour, cfg)

			return JSONSuccessWithExtra(c, "Login successful", nil, fiber.Map{"ok": true})
		}

		// For other purposes, just return success
		return JSONSuccessWithExtra(c, "OTP verified successfully", fiber.Map{"requestId": otpDoc.RequestID}, fiber.Map{"ok": true, "requestId": otpDoc.RequestID})
	})

	// POST /auth/otp/resend
	otpGroup.Post("/resend", func(c *fiber.Ctx) error {
		var body struct {
			Email   string `json:"email"`
			Purpose string `json:"purpose"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}

		// Validate email
		email := strings.TrimSpace(body.Email)
		if email == "" {
			return fiber.NewError(fiber.StatusBadRequest, "email is required")
		}

		// Validate purpose
		purpose := models.OTPPurpose(body.Purpose)
		if !purpose.IsValid() {
			return fiber.NewError(fiber.StatusBadRequest, "invalid purpose")
		}

		// Resend OTP
		requestID, err := otpService.ResendOTP(c.Context(), email, purpose)
		if err != nil {
			log.Error("OTP resend failed", zap.Error(err), zap.String("email", email))
			// Return generic message to avoid enumeration
			return JSONSuccessWithExtra(c, "If an account exists, we've sent a new code to your email", nil, fiber.Map{"ok": true})
		}

		return JSONSuccessWithExtra(c, "OTP resent", fiber.Map{"requestId": requestID}, fiber.Map{"requestId": requestID, "ok": true})
	})
}

func createJWT(secret string, u *models.User, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": u.Email,
		"exp":   time.Now().Add(ttl).Unix(),
		"iat":   time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func parseJWT(secret, token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, err
	}
	return parsed.Claims.(jwt.MapClaims), nil
}

func setCookie(c *fiber.Ctx, name, val string, ttl time.Duration, cfg config.Config) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    val,
		Expires:  time.Now().Add(ttl),
		HTTPOnly: true,
		Secure:   cfg.SecureCookies,
		SameSite: "Lax",
		Domain:   cfg.CookieDomain,
		Path:     "/",
	})
}

func clearCookie(c *fiber.Ctx, name string, cfg config.Config) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   cfg.SecureCookies,
		SameSite: "Lax",
		Domain:   cfg.CookieDomain,
		Path:     "/",
	})
}

// naive; replace with bcrypt util
func hashPassword(p string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hash)
}