package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"os"
	"raven-club/internal/pkg/user"
	"strconv"
)

var (
	ErrMissingFields      = errors.New("missing required fields")
	ErrPasswordProcess    = errors.New("failed to process password")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFailedUserCreation = errors.New("failed to create user")
)

type Service interface {
	Register(ctx *gin.Context, req RegisterRequest) (*Response, error)
	Login(ctx *gin.Context, req LoginRequest) (*Response, error)
}

type service struct {
	logger         *zap.Logger
	userService    user.Service
	emailLookup    *EmailLookup
	emailValidator *EmailValidator
	cookieStore    *sessions.CookieStore //TODO: move this to a module
}

func NewService(logger *zap.Logger, us user.Service) Service {
	return &service{
		logger:         logger, // CreateUser context to all logs
		userService:    us,
		emailLookup:    NewEmailLookup(us),
		emailValidator: NewEmailValidator(),
		cookieStore:    sessions.NewCookieStore([]byte(os.Getenv("TEST_SESSION_KEY"))), //TODO: move this to module
	}
}

// Register accepts RegisterRequest, creates a User, adds User to UserService, returns a JwtResponse
func (s *service) Register(ctx *gin.Context, req RegisterRequest) (*Response, error) {
	// CreateUser request logging
	s.logger.Info("processing registration request",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	// Generate user ID
	id, err := uuid.NewUUID()
	if err != nil {
		return &Response{}, err
	}

	// Validate and normalize email
	if err := s.emailValidator.Validate(req.Email); err != nil {
		s.logger.Warn("invalid email format",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return &Response{}, err
	}
	email := s.emailValidator.Normalize(req.Email)

	// Check if email exists
	if err := s.emailLookup.CheckEmailExists(ctx, email); err != nil && errors.Is(err, ErrEmailNotFound) {
		s.logger.Info("email does not exists",
			zap.String("email", req.Email),
		)
	} else if err != nil && errors.Is(err, ErrEmailExists) {
		s.logger.Error("Email exists", zap.Error(err))
		return &Response{}, err
	} else {
		s.logger.Error("failed to check email", zap.String("email", req.Email), zap.Error(err))
		return &Response{}, err
	}

	// Check required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		s.logger.Error("missing required fields",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		return &Response{}, ErrMissingFields
	}

	if req.Password != req.ConfirmPassword {
		s.logger.Warn("passwords do not match")
		return &Response{}, ErrPasswordProcess
	}

	// Hash password
	// TODO: salt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password",
			zap.Error(err),
		)
		return &Response{}, ErrPasswordProcess
	}

	u := &user.User{
		ID:       id.ID(),
		Username: req.Username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// CreateUser the user
	err = s.userService.CreateUser(ctx, *u)
	if err != nil {
		return &Response{}, ErrFailedUserCreation
	}

	// Create new user session
	var store = s.cookieStore
	userSession, err := store.Get(ctx.Request, strconv.Itoa(int(u.ID)))
	if err != nil {
		s.logger.Error("failed to get session", zap.Error(err))
		return &Response{}, err
	}

	s.logger.Info("user registered successfully",
		zap.String("username", u.Username),
		zap.String("email", u.Email),
		zap.Uint32("id", u.ID),
	)
	return &Response{
		Success: true,
		Message: "Registration successful",
		Session: userSession,
	}, nil
}

func (s *service) Login(ctx *gin.Context, req LoginRequest) (*Response, error) {
	// TODO implement me
	s.logger.Info("processing login request",
		zap.String("email", req.Email),
	)

	// Find and validate user
	u, err := s.emailLookup.FindUser(ctx, req.Email)
	if err != nil {
		s.logger.Info("login failed: user not found",
			zap.String("email", req.Email),
		)
		return &Response{}, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		s.logger.Info("login failed: invalid password",
			zap.String("email", req.Email),
		)
		return &Response{}, ErrInvalidCredentials
	}

	s.logger.Info("user logged in successfully",
		zap.String("email", req.Email),
		zap.Uint32("id", u.ID),
	)

	// Create new user session
	var store = s.cookieStore
	userSession, err := store.Get(ctx.Request, strconv.Itoa(int(u.ID)))
	// Add user to the session
	userSession.Values["username"] = u.Username
	userSession.Values["email"] = u.Email
	userSession.Values["Permissions"] = "Admin"

	if err != nil {
		s.logger.Error("failed to get session", zap.Error(err))
		return &Response{}, err
	}

	return &Response{
		Success: true,
		Message: "Login successful",
		Session: userSession,
	}, nil
}
