package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

// NewManager creates a new instance of Manager with the provided signingKey.
//
// Parameters:
//   - signingKey: A string used to sign tokens. Must not be empty.
//
// Returns:
//   - *Manager: A pointer to the newly created Manager instance.
//   - error: An error if the signingKey is empty.
func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

// NewJWT creates a new JWT token containing the provided user ID and TTL.
//
// The subject of the token will be set to the userId, and the expiration time
// will be set to the current time plus the provided ttl.
//
// Parameters:
//   - userId: The user ID to be included in the token.
//   - ttl: The TTL for which the token will remain valid.
//
// Returns:
//   - string: The signed JWT token.
//   - error: An error if the token could not be signed.
func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   userId,
	})

	return token.SignedString([]byte(m.signingKey))
}

// Parse verifies the provided accessToken and returns the user ID contained
// within the subject claim (sub) if the token is valid.
//
// Parameters:
//   - accessToken: The JWT token to be verified and parsed.
//
// Returns:
//   - string: The user ID contained in the subject claim of the token.
//   - error: An error if the token is invalid, or if the subject claim is not
//     present or could not be parsed as a string.
func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

// NewRefreshToken generates a cryptographically secure random string, which can be used to generate a refresh token.
//
// Returns:
//   - string: A random string suitable for use as a refresh token.
//   - error: An error if the random number generator fails.
func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
