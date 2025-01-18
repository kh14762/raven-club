package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"raven-club/internal/pkg/types"
	"raven-club/internal/pkg/user"
	"regexp"
	"strings"
)

// Email-related errors
var (
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrEmailExists   = errors.New("email already registered")
	ErrEmailNotFound = errors.New("email not found")
)

// EmailValidator handles email validation and lookup
type EmailValidator struct {
	emailRegex *regexp.Regexp
}

// NewEmailValidator creates a new email validator
func NewEmailValidator() *EmailValidator {
	return &EmailValidator{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
	}
}

// Validate checks if an email is valid
func (ev *EmailValidator) Validate(email string) error {
	if !ev.emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// Normalize converts email to consistent format (lowercase)
func (ev *EmailValidator) Normalize(email string) string {
	return strings.ToLower(email)
}

// EmailLookup handles email-based user lookups
type EmailLookup struct {
	userService user.Service
	validator   *EmailValidator
}

// NewEmailLookup creates a new email lookup service
func NewEmailLookup(us user.Service) *EmailLookup {
	return &EmailLookup{
		userService: us,
		validator:   NewEmailValidator(),
	}
}

// CheckEmailExists checks if an email is already registered
// returns true if email exists
func (el *EmailLookup) CheckEmailExists(ctx *gin.Context, email string) error {
	_, err := el.FindUser(ctx, email)
	if err == nil {
		return ErrEmailExists
	}
	if errors.Is(err, ErrEmailNotFound) {
		return err
	}
	return err
}

// FindUser looks up a user by email
func (el *EmailLookup) FindUser(ctx *gin.Context, email string) (*types.User, error) {
	// Validate email format
	if err := el.validator.Validate(email); err != nil {
		return &types.User{}, err
	}

	// Normalize email
	normalizedEmail := el.validator.Normalize(email)

	// Search for user
	users, err := el.userService.ListUsers(ctx)
	if err != nil {
		return &types.User{}, ErrEmailNotFound
	}
	for _, u := range users {
		if el.validator.Normalize(u.Email) == normalizedEmail {
			return &u, nil
		}
	}

	return &types.User{}, ErrEmailNotFound
}
