package models

// APIResponse represents the standard API response format
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	ErrorCode int         `json:"error_code,omitempty"`
	Data      interface{} `json:"data"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) APIResponse {
	if data == nil {
		data = map[string]interface{}{}
	}
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, errorCode int, data interface{}) APIResponse {
	if data == nil {
		data = map[string]interface{}{}
	}
	return APIResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Data:      data,
	}
}

// Error codes
const (
	// 10XX : Main App Errors
	ErrAppServer        = 1000 // Global Error
	ErrMissingHeaders   = 1001
	ErrMissingParams    = 1002
	ErrInvalidPagination = 1003
	ErrInvalidLocale    = 1004
	ErrInvalidTimezone  = 1005
	ErrRateLimited      = 1006

	// 11XX : Http Errors
	ErrUnauthorized     = 1101
	ErrForbidden        = 1102
	ErrUnprocessable    = 1103
	ErrAuthFailed       = 1104
	ErrNotFound         = 1105

	// 12XX : Auth Errors
	ErrTokenExpired     = 1201
	ErrInvalidSession   = 1202
	ErrInvalidToken     = 1204
	ErrUnauthenticated  = 1205
	ErrUserNotFound     = 1206

	// 13XX : Session Errors
	ErrInvalidCredentials = 1301
	ErrInvalidLoginType   = 1302
	ErrInvalidSocialType  = 1303
	ErrLoginError         = 1304
	ErrAccountDisabled    = 1305
	ErrInvalidMobile      = 1306
	ErrInvalidOTP         = 1307
	ErrInvalidEmailPass   = 1308
	ErrAccountExists      = 1309
	ErrRequestInvalid     = 1310
	ErrUnauthorizedAccess = 1311
	ErrActiveDirectory    = 1312
	ErrEmailNotConfirmed  = 1313
	ErrEmailLinkExpired   = 1314
	ErrAccountNotActive   = 1315
	ErrCannotDeleteUser   = 1316
	ErrMobileRegistered   = 1317
	ErrGoogleSignUp       = 1318
	ErrWrongOldMobile     = 1319
	ErrOTPExpired         = 1320
	ErrCannotDeleteProvider = 1321
	ErrAccountBlocked     = 1322
)
