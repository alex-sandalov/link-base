package domain

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
}
