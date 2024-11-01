package in_memory_redis

import (
	"context"
	"fmt"
	"link-base/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ReferralRedis struct {
	redisClient *redis.Client
}

// NewReferralRedis creates a new instance of ReferralRedis.
func NewReferralRedis(client *redis.Client) *ReferralRedis {
	return &ReferralRedis{
		redisClient: client,
	}
}

// Create sets a referral code in Redis with a TTL.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - referral: A domain.Referral struct containing the referral code and user ID.
//
// Returns:
//   - error: An error if the referral code can't be created in Redis.
func (r *ReferralRedis) Create(ctx context.Context, referral domain.Referral) error {
	fmt.Println("Creating referral code")

	_, err := r.redisClient.Set(ctx, referral.ReferralCode, referral.UserId.String(), referral.TTL).Result() // Приводим UserId к строке
	if err != nil {
		return fmt.Errorf("error setting referral code in Redis: %w", err)
	}

	return nil
}

// FindByReferralCode retrieves the creator of the referral code from Redis.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - referralCode: The referral code to search for.
//
// Returns:
//   - uuid.UUID: The user ID of the referral code creator if found.
//   - error: An error if the referral code can't be found in Redis.
func (r *ReferralRedis) FindByReferralCode(ctx context.Context, referralCode string) (uuid.UUID, error) {
	creatorIDStr, err := r.redisClient.Get(ctx, referralCode).Result()
	if err != nil {
		if err == redis.Nil {
			return uuid.Nil, fmt.Errorf("referral code not found: %s", referralCode)
		}
		return uuid.Nil, fmt.Errorf("error getting referral code from Redis: %w", err)
	}

	id, err := uuid.Parse(creatorIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing user ID from referral code: %w", err)
	}

	return id, nil
}
