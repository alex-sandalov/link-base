package postgres

import (
	"context"
	"fmt"
	"link-base/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReferralPostgres struct {
	db *sqlx.DB
}

// NewReferralPostgres creates a new instance of ReferralPostgres.
//
// Parameters:
//   - db: A pointer to a sqlx database connection.
//   - logger: A pointer to a slog logger.
//
// Returns:
//   - *ReferralPostgres: A new instance of ReferralPostgres.
func NewReferralPostgres(db *sqlx.DB) *ReferralPostgres {
	return &ReferralPostgres{
		db: db,
	}
}

// CreateReferral creates a new referral in the database.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - user: A domain.ReferralUser struct containing the user ID and referral ID.
//
// Returns:
//   - error: An error if the referral can't be created in the database.
//
// Note: The "ON CONFLICT (user_id) DO NOTHING" allows us to ignore the error if the user ID already exists in the
// table when inserting a new referral. This is useful when a user tries to refer someone who already has an account.
func (r *ReferralPostgres) CreateReferral(ctx context.Context, tx *sqlx.Tx, user domain.ReferralUser) error {
	const insertQuery = `
		INSERT INTO referral (user_id, referred_by_user_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING
	`

	_, err := tx.ExecContext(ctx, insertQuery, user.UserID, user.Referral)
	return err
}

// CreateReferralCode creates a new referral code in the database.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - referral: A domain.Referral struct containing the referral code and user ID.
//
// Returns:
//   - error: An error if the referral code can't be created in the database.
func (r *ReferralPostgres) CreateReferralCode(ctx context.Context, referral domain.Referral) error {
	const insertQuery = `
		INSERT INTO referral_code (user_id, code, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, code) DO UPDATE
		SET code = $2, expires_at = $3
	`

	ExpiresAt := time.Now().Add(referral.TTL)
	_, err := r.db.ExecContext(ctx, insertQuery, referral.UserId, referral.ReferralCode, ExpiresAt)
	if err != nil {
		return fmt.Errorf("error inserting or updating referral: %w", err)
	}

	return nil
}

// FindCodeByUserID retrieves all referrals associated with the given user ID from the database.
//
// The function executes a SQL query to select the user_id, code, and expires_at
// columns from the referrals table where the user_id matches the provided UUID.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - id: The UUID of the user whose referrals are to be retrieved.
//
// Returns:
//   - []domain.Referral: A slice of referral details if found.
//   - error: An error if there is a database query failure.
func (d *ReferralPostgres) FindCodeByUserID(ctx context.Context, id uuid.UUID) ([]domain.Referral, error) {
	var referrals []domain.Referral

	const findQuery = `
		SELECT user_id, code
		FROM referral_code
		WHERE user_id = $1 AND expires_at > NOW()
	`

	err := d.db.SelectContext(ctx, &referrals, findQuery, id)
	return referrals, err
}

// FindReferralByUserID retrieves all users that were referred by the given user ID.
//
// The function executes a SQL query to select the user_id
// column from the referral table where the referred_by_user_id matches the provided UUID.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - id: The UUID of the user whose referrals are to be retrieved.
//
// Returns:
//   - []domain.User: A slice of user details if found.
//   - error: An error if there is a database query failure.
func (d *ReferralPostgres) FindReferralByUserID(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	var users []uuid.UUID

	const findQuery = `
		SELECT user_id
		FROM referral
		WHERE referred_by_user_id = $1
	`

	err := d.db.SelectContext(ctx, &users, findQuery, id)
	return users, err
}
