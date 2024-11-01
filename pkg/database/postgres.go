package database

import (
	"fmt"
	"link-base/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresClient initializes and returns a connection to a PostgreSQL database using the provided configuration.
//
// Parameters:
//   - cfg: A PostgresConfig struct containing the database connection details such as host, port, user, database name, password, and SSL mode.
//
// Returns:
//   - *sqlx.DB: A pointer to the initialized database connection.
//   - error: An error if the connection to the database fails.
func NewPostgresClient(cfg config.PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Database, cfg.Password, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
