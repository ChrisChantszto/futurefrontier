package models

import "time"

// OTPPurpose represents the purpose of an OTP code
type OTPPurpose string

const (
	OTPPurposeLogin          OTPPurpose = "login"
	OTPPurposeForgetPassword OTPPurpose = "forgetPassword"
)

// IsValid checks if the OTP purpose is valid
func (p OTPPurpose) IsValid() bool {
	switch p {
	case OTPPurposeLogin, OTPPurposeForgetPassword:
		return true
	default:
		return false
	}
}

// OTPCode represents an OTP code document in MongoDB
type OTPCode struct {
	ID          string     `bson:"_id,omitempty" json:"id,omitempty"`
	Email       string     `bson:"email" json:"email"`
	RequestID   string     `bson:"requestId" json:"requestId"`
	Purpose     OTPPurpose `bson:"purpose" json:"purpose"`
	CodeHMAC    string     `bson:"codeHmac" json:"-"` // Don't expose in JSON
	ExpiresAt   time.Time  `bson:"expiresAt" json:"expiresAt"`
	Attempts    int        `bson:"attempts" json:"attempts"`
	MaxAttempts int        `bson:"maxAttempts" json:"maxAttempts"`
	CreatedAt   time.Time  `bson:"createdAt" json:"createdAt"`
	ConsumedAt  *time.Time `bson:"consumedAt" json:"consumedAt"`
	LastSentAt  time.Time  `bson:"lastSentAt" json:"lastSentAt"`
}

// IsExpired checks if the OTP code has expired
func (o *OTPCode) IsExpired() bool {
	return time.Now().UTC().After(o.ExpiresAt)
}

// IsConsumed checks if the OTP code has been consumed
func (o *OTPCode) IsConsumed() bool {
	return o.ConsumedAt != nil
}

// CanAttempt checks if more attempts are allowed
func (o *OTPCode) CanAttempt() bool {
	return o.Attempts < o.MaxAttempts
}

// IsActive checks if the OTP is active (not expired, not consumed, can attempt)
func (o *OTPCode) IsActive() bool {
	return !o.IsExpired() && !o.IsConsumed() && o.CanAttempt()
}
