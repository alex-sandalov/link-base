package domain

import "github.com/google/uuid"

type ReferralUser struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Referral uuid.UUID `json:"referral" db:"referral"`
}
