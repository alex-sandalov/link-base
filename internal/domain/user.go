package domain

import "github.com/google/uuid"

type User struct {
	UserId       uuid.UUID `db:"user_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
}
