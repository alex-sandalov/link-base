package hash

import (
	"crypto/sha1"
	"fmt"
)

// PasswordHasher provides hashing logic to securely store passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
}

// SHA1Hasher uses SHA1 to hash passwords with provided salt.
type SHA1Hasher struct {
	salt string
}

// NewSHA1Hasher creates a new SHA1Hasher instance with the provided salt.
//
// Parameters:
//   - salt: A string used to salt the hash, enhancing its security.
//
// Returns:
//   - *SHA1Hasher: A pointer to the newly created SHA1Hasher instance.
func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

// Hash takes a password and returns a hashed version of it with the provided salt.
//
// The method uses SHA1 to hash the password with the provided salt.
//
// Parameters:
//   - password: The string to be hashed.
//
// Returns:
//   - string: The hashed password.
//   - error: An error if there was a problem while hashing the password.
func (h *SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
