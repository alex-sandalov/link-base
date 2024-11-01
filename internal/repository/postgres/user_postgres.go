package postgres

import (
	"context"
	"fmt"
	"link-base/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

// NewUserPostgres creates a new instance of UserPostgres.
//
// Parameters:
//   - db: A pointer to a sqlx database connection.
//   - logger: A pointer to a slog logger.
//
// Returns:
//   - *UserPostgres: A new instance of UserPostgres.
func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

// Create creates a new user in the database.
//
// The method executes a SQL query to insert a new user into the users table.
// The context is used to pass request-scoped values to the database driver.
//
// The method returns an error if the user already exists in the database.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - tx: A pointer to a sqlx transaction.
//   - u: The user to be created, containing the user ID, email, and password hash.
//
// Returns:
//   - error: An error if the user already exists in the database.
func (d *UserPostgres) Create(ctx context.Context, tx *sqlx.Tx, u domain.User) error {
	const queryCreate = `
		INSERT INTO users (user_id, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
	`
	_, err := tx.ExecContext(ctx, queryCreate, u.UserId, u.Email, u.PasswordHash)
	return err
}

// FindByUserId retrieves a user from the database by their unique user ID.
//
// The function executes a SQL query to select the user_id, email, and password_hash
// columns from the users table where the user_id matches the provided UUID.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - userId: The UUID of the user to be retrieved.
//
// Returns:
//   - domain.User: The user details if found.
//   - error: An error if the user is not found or if there is a database query failure.
func (d *UserPostgres) FindByUserId(ctx context.Context, userId uuid.UUID) (domain.User, error) {
	var usr domain.User
	const findQuery = `
		SELECT user_id, email, password_hash
		FROM users
		WHERE user_id = $1
		LIMIT 1
	`

	if err := d.db.GetContext(ctx, &usr, findQuery, userId); err != nil {
		return usr, fmt.Errorf("could not find user with ID %s: %w", userId, err)
	}

	return usr, nil
}

// FindByEmail retrieves a user from the database by their unique email address.
//
// The function executes a SQL query to select the user_id, email, and password_hash
// columns from the users table where the email matches the provided string.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - email: The email address of the user to be retrieved.
//
// Returns:
//   - domain.User: The user details if found.
//   - error: An error if the user is not found or if there is a database query failure.
func (d *UserPostgres) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	const findQuery = `
		SELECT user_id, email, password_hash
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var user domain.User
	if err := d.db.GetContext(ctx, &user, findQuery, email); err != nil {
		return domain.User{
			UserId: uuid.Nil,
		}, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}
