package domain

import (
	"github.com/google/uuid"
	"time"
)

type Referral struct {
	ReferralCode string        `db:"code"`
	UserId       uuid.UUID     `db:"user_id"`
	TTL          time.Duration `db:"ttl"`
}
