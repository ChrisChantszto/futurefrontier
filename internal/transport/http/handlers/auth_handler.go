package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			return JSONErrorAlways200(c, "invalid body", 1001, nil)
		}

		ctx := c.Context()
		count, err := svc.CountUsers(ctx)
		if err != nil {
			return JSONErrorAlways200(c, "db error", 1002, nil)
		}
		if count > 0 && body.Email != cfg.SuperAdminEmail {
			return JSONErrorAlways200(c, "already initialized", 1003, nil)
		}
		_, err = svc.CreateUser(ctx, body.Email, body.Password, []models.Role{models.RoleSuperAdmin})
		if err != nil {
			return JSONErrorAlways200(c, "create user failed", 1004, nil)
		}
		return JSONSuccessWithExtra(c, "User initialized", nil, fiber.Map{"ok": true})
	})

	group.Post("/login", func(c *fiber.Ctx) error {
		var body struct{ Email, Password string }
		if err := c.BodyParser(&body); err != nil {
			return JSONErrorAlways200(c, "invalid body", 3001, nil)
		}
		u, err := svc.FindUserByEmail(c.Context(), body.Email)
		if err != nil || !svc.VerifyPassword(u, body.Password) {
			return JSONErrorAlways200(c, "invalid credentials", 3002, nil)
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
			return JSONErrorAlways200(c, "no refresh token", 4001, nil)
		}
		claims, err := parseJWT(cfg.JWTRefreshSecret, rt)
		if err != nil {
			return JSONErrorAlways200(c, "invalid refresh token", 4002, nil)
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
			return JSONErrorAlways200(c, "invalid body", 5001, nil)
		}
		rt, err := svc.ConsumeResetToken(c.Context(), body.Token)
		if err != nil {
			return JSONErrorAlways200(c, "invalid token", 5002, nil)
		}
		// update user password by ObjectID
		oid, err := primitive.ObjectIDFromHex(rt.UserID)
		if err != nil {
			return JSONErrorAlways200(c, "invalid user id", 5003, nil)
		}
		_, err = db.Collection("users").UpdateByID(c.Context(), oid, bson.M{"$set": bson.M{
			"passwordHash": hashPassword(body.Password),
			"updatedAt":   time.Now().UTC(),
		}})
		if err != nil {
			return JSONErrorAlways200(c, "update failed", 5004, nil)
		}
		return JSONSuccessWithExtra(c, "Password reset successfully", nil, fiber.Map{"ok": true})
	})

	// TOTP endpoints (email-based one-time codes with time limit)
	// POST /auth/totp/request - request TOTP for login (default) or forgetPassword purpose
	group.Post("/totp/request", func(c *fiber.Ctx) error {
		var body struct {
			Email   string `json:"email"`
			Purpose string `json:"purpose"`
		}
		if err := c.BodyParser(&body); err != nil {
			return JSONErrorAlways200(c, "invalid body", 7001, nil)
		}
		email := strings.TrimSpace(body.Email)
		if email == "" {
			return JSONErrorAlways200(c, "email is required", 7002, nil)
		}
		purpose := models.OTPPurpose(body.Purpose)
		if body.Purpose == "" { // default to login
			purpose = models.OTPPurposeLogin
		}
		if !purpose.IsValid() {
			return JSONErrorAlways200(c, "invalid purpose", 7003, nil)
		}
		requestID, err := otpService.RequestOTP(c.Context(), email, purpose)
		if err != nil {
			log.Error("TOTP request failed", zap.Error(err), zap.String("email", email))
			return JSONSuccessWithExtra(c, "If an account exists, we've sent a code to your email", nil, nil)
		}
		return JSONSuccessWithExtra(c, "TOTP sent", fiber.Map{"requestId": requestID}, fiber.Map{"requestId": requestID})
	})

	// POST /auth/totp/login - verify TOTP for login and issue JWT cookies
	group.Post("/totp/login", func(c *fiber.Ctx) error {
		var body struct {
			Email     string `json:"email"`
			Code      string `json:"code"`
			RequestID string `json:"requestId,omitempty"`
		}
		if err := c.BodyParser(&body); err != nil {
			return JSONErrorAlways200(c, "invalid body", 8001, nil)
		}
		email := strings.TrimSpace(body.Email)
		code := strings.TrimSpace(body.Code)
		if email == "" || code == "" {
			return JSONErrorAlways200(c, "email and code are required", 8002, nil)
		}
		_, err := otpService.VerifyOTP(c.Context(), email, code, body.RequestID, models.OTPPurposeLogin)
		if err != nil {
			log.Error("TOTP login verification failed", zap.Error(err), zap.String("email", email))
			return JSONErrorAlways200(c, "invalid or expired code", 8003, nil)
		}
		// Find user and issue tokens
		user, err := svc.FindUserByEmail(c.Context(), email)
		if err != nil {
			return JSONErrorAlways200(c, "user not found", 8004, nil)
		}
		accessToken, _ := createJWT(cfg.JWTSecret, user, 15*time.Minute)
		refreshToken, _ := createJWT(cfg.JWTRefreshSecret, user, 7*24*time.Hour)
		setCookie(c, "access", accessToken, 15*time.Minute, cfg)
		setCookie(c, "refresh", refreshToken, 7*24*time.Hour, cfg)
		return JSONSuccessWithExtra(c, "Login successful", nil, fiber.Map{"ok": true})
	})

	// Forgot password via TOTP: request and verify
	// POST /auth/forgot-password/request-totp
	group.Post("/forgot-password/request-totp", func(c *fiber.Ctx) error {
		var body struct{ Email string `json:"email"` }
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		email := strings.TrimSpace(body.Email)
		if email == "" {
			return fiber.NewError(fiber.StatusBadRequest, "email is required")
		}
		requestID, err := otpService.RequestOTP(c.Context(), email, models.OTPPurposeForgetPassword)
		if err != nil {
			log.Error("Forgot-password TOTP request failed", zap.Error(err), zap.String("email", email))
			return JSONSuccessWithExtra(c, "If an account exists, we've sent a code to your email", nil, nil)
		}
		return JSONSuccessWithExtra(c, "TOTP sent", fiber.Map{"requestId": requestID}, fiber.Map{"requestId": requestID})
	})

	// POST /auth/forgot-password/verify-totp
	group.Post("/forgot-password/verify-totp", func(c *fiber.Ctx) error {
		var body struct {
			Email     string `json:"email"`
			Code      string `json:"code"`
			RequestID string `json:"requestId,omitempty"`
			Password  string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return JSONErrorAlways200(c, "invalid body", 9001, nil)
		}
		email := strings.TrimSpace(body.Email)
		code := strings.TrimSpace(body.Code)
		if email == "" || code == "" || strings.TrimSpace(body.Password) == "" {
			return JSONErrorAlways200(c, "email, code and password are required", 9002, nil)
		}
		_, err := otpService.VerifyOTP(c.Context(), email, code, body.RequestID, models.OTPPurposeForgetPassword)
		if err != nil {
			log.Error("Forgot-password TOTP verification failed", zap.Error(err), zap.String("email", email))
			return JSONErrorAlways200(c, "invalid or expired code", 9003, nil)
		}
		// Update password for the user
		u, err := svc.FindUserByEmail(c.Context(), email)
		if err != nil {
			return JSONErrorAlways200(c, "user not found", 9004, nil)
		}
		// Update by email to avoid ObjectID vs string mismatch on _id
		_, err = db.Collection("users").UpdateOne(c.Context(), bson.M{"email": u.Email}, bson.M{"$set": bson.M{
			"passwordHash": hashPassword(body.Password),
			"updatedAt":   time.Now().UTC(),
		}})
		if err != nil {
			return JSONErrorAlways200(c, "update failed", 9005, nil)
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
			return JSONErrorAlways200(c, "invalid body", 2001, nil)
		}

		// Validate inputs
		email := strings.TrimSpace(body.Email)
		code := strings.TrimSpace(body.Code)
		if email == "" || code == "" {
			return JSONErrorAlways200(c, "email and code are required", 2002, nil)
		}

		// Validate purpose
		purpose := models.OTPPurpose(body.Purpose)
		if !purpose.IsValid() {
			return JSONErrorAlways200(c, "invalid purpose", 2003, nil)
		}

		// Verify OTP
		otpDoc, err := otpService.VerifyOTP(c.Context(), email, code, body.RequestID, purpose)
		if err != nil {
			log.Error("OTP verification failed", zap.Error(err), zap.String("email", email))
			return JSONErrorAlways200(c, "invalid or expired code", 2004, nil)
		}

		// For login purpose, create JWT tokens
		if purpose == models.OTPPurposeLogin {
			// Find user by email
			user, err := svc.FindUserByEmail(c.Context(), email)
			if err != nil {
				log.Error("User not found after OTP verification", zap.Error(err), zap.String("email", email))
				return JSONErrorAlways200(c, "user not found", 2005, nil)
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
			return JSONErrorAlways200(c, "invalid body", 6001, nil)
		}

		// Validate email
		email := strings.TrimSpace(body.Email)
		if email == "" {
			return JSONErrorAlways200(c, "email is required", 6002, nil)
		}

		// Validate purpose
		purpose := models.OTPPurpose(body.Purpose)
		if !purpose.IsValid() {
			return JSONErrorAlways200(c, "invalid purpose", 6003, nil)
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