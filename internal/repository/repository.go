package repository

import (
	"context"
	"link-base/internal/domain"
	"link-base/internal/repository/postgres"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
)

type User interface {
	Create(ctx context.Context, tx *sqlx.Tx, user domain.User) error
	FindByUserId(ctx context.Context, id uuid.UUID) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type RefreshToken interface {
	Create(ctx context.Context, session domain.RefreshToken) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (domain.RefreshToken, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error)
}

type Referral interface {
	CreateReferral(ctx context.Context, tx *sqlx.Tx, user domain.ReferralUser) error
	FindReferralByUserID(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
	CreateReferralCode(ctx context.Context, referral domain.Referral) error
	FindCodeByUserID(ctx context.Context, id uuid.UUID) ([]domain.Referral, error)
}

type Repository struct {
	User         User
	RefreshToken RefreshToken
	Referral     Referral
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:         postgres.NewUserPostgres(db),
		RefreshToken: postgres.NewRefreshTokenPostgres(db),
		Referral:     postgres.NewReferralPostgres(db),
	}
}
