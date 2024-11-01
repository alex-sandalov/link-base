package postgres

import (
	"context"
	"fmt"
	"link-base/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenPostgres struct {
	db *sqlx.DB
}

// NewRefreshTokenPostgres creates a new instance of RefreshTokenPostgres.
//
// Parameters:
//   - db: A pointer to a sqlx database connection.
//   - logger: A pointer to a slog logger.
//
// Returns:
//   - *RefreshTokenPostgres: A new instance of RefreshTokenPostgres.
func NewRefreshTokenPostgres(db *sqlx.DB) *RefreshTokenPostgres {
	return &RefreshTokenPostgres{
		db: db,
	}
}

// Create inserts a new refresh token into the database.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - refreshToken: The refresh token to be inserted, including user ID and expiration.
//
// Returns:
//   - error: An error if the insertion or update fails.
func (r *RefreshTokenPostgres) Create(ctx context.Context, refreshToken domain.RefreshToken) error {
	const insertQuery = `
		INSERT INTO refresh_token (user_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, refresh_token) DO UPDATE
		SET refresh_token = $2, expires_at = $3
	`

	_, err := r.db.ExecContext(ctx, insertQuery, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt)
	if err != nil {
		return fmt.Errorf("error inserting or updating refresh token: %w", err)
	}

	return nil
}

// DeleteByUserID deletes all refresh tokens associated with the given user ID from the database.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - userID: The UUID of the user whose refresh tokens are to be deleted.
//
// Returns:
//   - error: An error if the deletion fails.
func (r *RefreshTokenPostgres) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	const deleteQuery = `
		DELETE FROM refresh_token
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, deleteQuery, userID)
	return err
}

// FindByUserID retrieves a refresh token from the database by the user's unique user ID.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - userID: The UUID of the user whose refresh token is to be retrieved.
//
// Returns:
//   - domain.RefreshToken: The refresh token details if found.
//   - error: An error if the refresh token is not found or if there is a database query failure.
func (r *RefreshTokenPostgres) FindByUserID(ctx context.Context, userID uuid.UUID) (domain.RefreshToken, error) {
	const findQuery = `
		SELECT user_id, refresh_token, expires_at
		FROM refresh_token
		WHERE user_id = $1 AND expires_at > NOW()
	`

	var refreshToken domain.RefreshToken
	if err := r.db.GetContext(ctx, &refreshToken, findQuery, userID); err != nil {
		return domain.RefreshToken{}, fmt.Errorf("refresh token not found for user ID %s: %w", userID, err)
	}

	return refreshToken, nil
}

// FindByRefreshToken retrieves a refresh token from the database by the refresh token itself.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - refreshToken: The refresh token to be retrieved.
//
// Returns:
//   - domain.RefreshToken: The refresh token details if found.
//   - error: An error if the refresh token is not found or if there is a database query failure.
func (r *RefreshTokenPostgres) FindByRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error) {
	const findQuery = `
		SELECT user_id, refresh_token, expires_at
		FROM refresh_token
		WHERE refresh_token = $1 AND expires_at > NOW()
		LIMIT 1
	`

	var refreshTokenFromDB domain.RefreshToken
	err := r.db.GetContext(ctx, &refreshTokenFromDB, findQuery, refreshToken)

	if err != nil {
		return domain.RefreshToken{}, fmt.Errorf("refresh token not found: %w", err)
	}

	return refreshTokenFromDB, nil
}
