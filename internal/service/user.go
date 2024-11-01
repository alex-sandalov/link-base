package service

import (
	"context"
	"fmt"
	"link-base/internal/cache"
	"link-base/internal/config"
	"link-base/internal/domain"
	"link-base/internal/repository"
	"link-base/pkg/auth"
	"link-base/pkg/hash"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
)

type CreateUserInput struct {
	Email      string
	Password   string
	ReferralId uuid.UUID
}

type UserService struct {
	db           *sqlx.DB
	repos        *repository.Repository
	logger       *slog.Logger
	cfg          config.JWTConfig
	tokenManager *auth.Manager
	hasher       *hash.SHA1Hasher
	redis        *cache.Cache
}

// NewUserService creates a new instance of UserService.
//
// Parameters:
//   - repos: A pointer to a repository.Repository instance.
//   - logger: A pointer to a slog.Logger instance for logging.
//   - cfg: The configuration settings for JWT tokens.
//   - tokenManager: A pointer to an auth.Manager for managing authentication tokens.
//   - hasher: A pointer to a hash.SHA1Hasher for password hashing.
//   - db: A pointer to a sqlx.DB instance for database interactions.
//
// Returns:
//   - *UserService: A new instance of UserService.
func NewUserService(repos *repository.Repository, logger *slog.Logger,
	cfg config.JWTConfig, tokenManager *auth.Manager, hasher *hash.SHA1Hasher, db *sqlx.DB, redis *cache.Cache) *UserService {
	return &UserService{
		repos:        repos,
		logger:       logger,
		cfg:          cfg,
		tokenManager: tokenManager,
		hasher:       hasher,
		db:           db,
		redis:        redis,
	}
}

// SignIn authenticates a user with the provided credentials and returns a new
// session (access and refresh tokens) if the authentication is successful.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - input: The SignInInput containing the email and password of the user to be
//     authenticated.
//
// Returns:
//   - Tokens: A Tokens object containing the access and refresh tokens for the
//     newly created session.
//   - error: An error if the authentication fails or if there is a database query
//     failure.
func (u *UserService) SignIn(ctx context.Context, input SignInInput) (Tokens, error) {
	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := u.repos.User.FindByEmail(ctx, input.Email)
	if err != nil {
		return Tokens{}, err
	}

	if user.PasswordHash != passwordHash {
		return Tokens{}, fmt.Errorf("invalid credentials")
	}

	return u.createSession(ctx, user.UserId)
}

// SignUp registers a new user with the provided credentials and returns a new session.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - input: The SignUpInput containing the email, password, and referral code for the user to be registered.
//
// Returns:
//   - Tokens: A Tokens object containing the access and refresh tokens for the newly created session.
//   - error: An error if registration fails or if there is a database query failure.
func (u *UserService) SignUp(ctx context.Context, input SignUpInput) (Tokens, error) {
	referralId := uuid.Nil
	if input.ReferralCode != "" {
		var err error
		referralId, err = u.redis.Referral.FindByReferralCode(ctx, input.ReferralCode)
		if err != nil {
			return Tokens{}, fmt.Errorf("failed to find referral code: %w", err)
		}
	}

	return u.createUser(ctx, CreateUserInput{
		Email:      input.Email,
		Password:   input.Password,
		ReferralId: referralId,
	})
}

// RefreshTokens generates a new set of tokens using the provided refresh token.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - refreshToken: The refresh token used to generate new session tokens.
//
// Returns:
//   - Tokens: A new set of access and refresh tokens.
//   - error: An error if the refresh token is invalid or if there is a database query failure.
func (u *UserService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	session, err := u.repos.RefreshToken.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, fmt.Errorf("failed to find refresh token: %w", err)
	}

	return u.createSession(ctx, session.UserID)
}

// createSession creates a new session for the given user ID and returns the session tokens.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - userID: The UUID of the user for whom the session is to be created.
//
// Returns:
//   - Tokens: The session tokens containing the access token and refresh token.
//   - error: An error if the session could not be created or if there is a database query failure.
func (u *UserService) createSession(ctx context.Context, userID uuid.UUID) (Tokens, error) {
	accessToken, err := u.tokenManager.NewJWT(userID.String(), u.cfg.AccessTokenTTL)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := u.tokenManager.NewRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	session := domain.RefreshToken{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.cfg.RefreshTokenTTL),
	}

	if err := u.repos.RefreshToken.Create(ctx, session); err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// createUser registers a new user with the provided email and password and returns a new session.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - input: The createUserInput containing the email, password, and referral ID for the user to be registered.
//
// Returns:
//   - Tokens: The session tokens containing the access token and refresh token.
//   - error: An error if the session could not be created or if there is a database query failure.
func (u *UserService) createUser(ctx context.Context, input CreateUserInput) (Tokens, error) {
	_, err := u.repos.User.FindByEmail(ctx, input.Email)
	if err == nil {
		return Tokens{}, fmt.Errorf("email already in use")
	}

	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return Tokens{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	user := domain.User{
		UserId:       uuid.New(),
		Email:        input.Email,
		PasswordHash: passwordHash,
	}

	if err = u.repos.User.Create(ctx, tx, user); err != nil {
		return Tokens{}, err
	}

	if input.ReferralId != uuid.Nil {
		if err = u.repos.Referral.CreateReferral(ctx, tx, domain.ReferralUser{
			UserID:   user.UserId,
			Referral: input.ReferralId,
		}); err != nil {
			return Tokens{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return Tokens{}, err
	}

	u.logger.Info("Create user")
	return u.createSession(ctx, user.UserId)
}
