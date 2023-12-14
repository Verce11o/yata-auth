package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserID       uuid.UUID `json:"id" db:"user_id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash []byte    `json:"passwordHash" db:"password"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type SendUserEmailRequest struct {
	To   string `json:"to"`
	Code string `json:"code"`
}

type VerificationCode struct {
	Code       uuid.UUID `json:"code" db:"code"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	ExpireDate time.Time `json:"expire_date" db:"expire_date"`
}
