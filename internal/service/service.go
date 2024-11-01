package service

import (
	"context"
	"link-base/internal/cache"
	"link-base/internal/config"
	"link-base/internal/repository"
	"link-base/pkg/auth"
	"link-base/pkg/hash"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type SignInInput struct {
	Email    string
	Password string
}

type SignUpInput struct {
	Email        string
	Password     string
	ReferralCode string
}

type ReferralInput struct {
	UserId uuid.UUID
	TTL    time.Duration
}

type User interface {
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)
	SignUp(ctx context.Context, input SignUpInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Referral interface {
	CreateCode(ctx context.Context, input ReferralInput) (string, error)
	FindReferralByUserID(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
	SendEmail(ctx context.Context, userId uuid.UUID, email string) error
}

type Service struct {
	User     User
	Referral Referral
}

func NewService(repos *repository.Repository, logger *slog.Logger,
	cfg config.JWTConfig, tokenManager *auth.Manager, hasher *hash.SHA1Hasher, db *sqlx.DB, redis *cache.Cache) *Service {
	return &Service{
		User:     NewUserService(repos, logger, cfg, tokenManager, hasher, db, redis),
		Referral: NewReferralService(repos, redis, tokenManager),
	}
}
