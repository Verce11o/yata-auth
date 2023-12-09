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
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
