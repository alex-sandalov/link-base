package service

import (
	"context"
	"fmt"
	"link-base/internal/cache"
	"link-base/internal/config"
	"link-base/internal/domain"
	"link-base/internal/repository"
	"link-base/pkg/auth"
	"net/smtp"

	"github.com/google/uuid"
)

type ReferralService struct {
	repos        *repository.Repository
	redis        *cache.Cache
	tokenManager *auth.Manager
	cfg          config.SMPTConfig
}

// NewReferralService creates a new instance of ReferralService.
//
// Parameters:
//   - repos: A pointer to a Repository instance.
//   - redis: A pointer to a Cache instance.
//   - tokenManager: A pointer to a Manager instance.
//
// Returns:
//   - *ReferralService: A new instance of ReferralService.
func NewReferralService(repos *repository.Repository, redis *cache.Cache, tokenManager *auth.Manager) *ReferralService {
	return &ReferralService{
		repos:        repos,
		redis:        redis,
		tokenManager: tokenManager,
	}
}

// CreateCode creates a new referral code with the given user ID and TTL.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - input: A ReferralInput struct containing the user ID and TTL.
//
// Returns:
//   - string: The referral code if created successfully.
//   - error: An error if the referral code can't be created.
func (r *ReferralService) CreateCode(ctx context.Context, input ReferralInput) (string, error) {
	referralCode, err := r.generateReferralCode()
	if err != nil {
		return "", err
	}

	res, err := r.repos.Referral.FindCodeByUserID(ctx, input.UserId)
	if res != nil {
		return "", fmt.Errorf("referral code %s already exists", res[0].ReferralCode)
	}
	if err != nil {
		return "", err
	}

	referral := domain.Referral{
		ReferralCode: referralCode,
		UserId:       input.UserId,
		TTL:          input.TTL,
	}

	if err = r.repos.Referral.CreateReferralCode(ctx, referral); err != nil {
		return "", err
	}

	if err = r.redis.Referral.Create(ctx, referral); err != nil {
		return "", err
	}

	return referralCode, nil
}

// FindReferralByUserID retrieves all referral user IDs associated with the given user ID.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - id: The UUID of the user whose referral IDs are to be retrieved.
//
// Returns:
//   - []uuid.UUID: A slice of referral user IDs if found.
//   - error: An error if there is a database query failure.
func (r *ReferralService) FindReferralByUserID(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	return r.repos.Referral.FindReferralByUserID(ctx, id)
}

// generateReferralCode generates a new cryptographically secure referral code.
//
// Returns:
//   - string: The generated referral code.
//   - error: An error if the code generation fails.
func (r *ReferralService) generateReferralCode() (string, error) {
	return r.tokenManager.NewRefreshToken()
}

// SendEmail sends an email containing the referral code to the specified email address.
//
// Parameters:
//   - ctx: The context for controlling the request lifecycle.
//   - userId: The UUID of the user whose referral code is to be sent.
//   - email: The recipient's email address.
//
// Returns:
//   - error: An error if sending the email fails.
func (r *ReferralService) SendEmail(ctx context.Context, userId uuid.UUID, email string) error {
	referralCode, err := r.repos.Referral.FindCodeByUserID(ctx, userId)
	if err != nil {
		return err
	}

	from := userId.String()
	subject := "Your Referral Code"
	body := fmt.Sprintf("Hello!\n\nYour referral code is: %s\n\nBest regards!", referralCode)
	message := []byte("Subject: " + subject + "\n\n" + body)

	smtpAuth := smtp.PlainAuth("", r.cfg.SMPTUser, r.cfg.SMPTPassword, r.cfg.SMPTHost)

	err = smtp.SendMail(r.cfg.SMPTHost+":"+r.cfg.SMPTPort, smtpAuth, from, []string{email}, message)
	if err != nil {
		return err
	}

	return nil
}
