package models

import "time"

type ResetToken struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"userId"`
	Hash      string    `bson:"hash"`
	ExpiresAt time.Time `bson:"expiresAt"`
	Used      bool      `bson:"used"`
}