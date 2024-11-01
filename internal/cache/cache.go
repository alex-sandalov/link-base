package cache

import (
	"context"
	InMemoryRedis "link-base/internal/cache/in-memory-redis"
	"link-base/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Referral interface {
	Create(ctx context.Context, referral domain.Referral) error
	FindByReferralCode(ctx context.Context, referralCode string) (uuid.UUID, error)
}

type Cache struct {
	Referral Referral
}

// NewCache initializes and returns a new Cache instance.
//
// Parameters:
//   - redisClient: A pointer to a Redis client used to interact with the Redis database.
//   - logger: A pointer to a slog logger for logging purposes.
//
// Returns:
//   - *Cache: A new instance of Cache.
func NewCache(redisClient *redis.Client) *Cache {
	return &Cache{
		Referral: InMemoryRedis.NewReferralRedis(redisClient),
	}
}
