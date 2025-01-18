package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"os"
	"raven-club/internal/pkg/session"
	"raven-club/internal/pkg/types"
	"raven-club/internal/pkg/user"
)

var (
	ErrMissingFields      = errors.New("missing required fields")
	ErrPasswordProcess    = errors.New("failed to process password")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFailedUserCreation = errors.New("failed to create user")
)

type Service interface {
	Register(ctx *gin.Context, req types.RegisterRequest) (*types.RegisterResponse, error)
	Login(ctx *gin.Context, req types.LoginRequest) (*types.LoginResponse, error)
}

type service struct {
	logger         *zap.Logger
	userService    user.Service
	sessionService session.Service
	emailLookup    *EmailLookup
	emailValidator *EmailValidator
	cookieStore    *sessions.CookieStore
}

func NewService(logger *zap.Logger, us user.Service, ss session.Service) Service {
	return &service{
		logger:         logger, // CreateUser context to all logs
		userService:    us,
		sessionService: ss,
		emailLookup:    NewEmailLookup(us),
		emailValidator: NewEmailValidator(),
		cookieStore:    sessions.NewCookieStore([]byte(os.Getenv("TEST_SESSION_KEY"))),
	}
}

// Register accepts RegisterRequest, creates a User, adds User to UserService, returns a JwtResponse
func (s *service) Register(ctx *gin.Context, req types.RegisterRequest) (*types.RegisterResponse, error) {
	// CreateUser request logging
	s.logger.Info("processing registration request",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	// Generate user ID
	id, err := uuid.NewUUID()
	if err != nil {
		return &types.RegisterResponse{}, err
	}

	// Validate and normalize email
	if err := s.emailValidator.Validate(req.Email); err != nil {
		s.logger.Warn("invalid email format",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return &types.RegisterResponse{}, err
	}
	email := s.emailValidator.Normalize(req.Email)

	// Check if email exists
	if err := s.emailLookup.CheckEmailExists(ctx, email); err != nil && errors.Is(err, ErrEmailNotFound) {
		s.logger.Info("email does not exists",
			zap.String("email", req.Email),
		)
	} else if err != nil && errors.Is(err, ErrEmailExists) {
		s.logger.Error("Email exists", zap.Error(err))
		return &types.RegisterResponse{}, err
	} else {
		s.logger.Error("failed to check email", zap.String("email", req.Email), zap.Error(err))
		return &types.RegisterResponse{}, err
	}

	// Check required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		s.logger.Error("missing required fields",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		return &types.RegisterResponse{}, ErrMissingFields
	}

	if req.Password != req.ConfirmPassword {
		s.logger.Warn("passwords do not match")
		return &types.RegisterResponse{}, ErrPasswordProcess
	}

	// Hash password
	// TODO: salt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password",
			zap.Error(err),
		)
		return &types.RegisterResponse{}, ErrPasswordProcess
	}

	u := &types.User{
		ID:       id.ID(),
		Username: req.Username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// CreateUser the user
	err = s.userService.CreateUser(ctx, *u)
	if err != nil {
		return &types.RegisterResponse{}, ErrFailedUserCreation
	}

	s.logger.Info("user registered successfully",
		zap.String("username", u.Username),
		zap.String("email", u.Email),
		zap.Uint32("id", u.ID),
	)
	return &types.RegisterResponse{
		Success: true,
		Message: "Registration successful",
		User:    u,
	}, nil
}

func (s *service) Login(ctx *gin.Context, req types.LoginRequest) (*types.LoginResponse, error) {
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
		return &types.LoginResponse{}, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		s.logger.Info("login failed: invalid password",
			zap.String("email", req.Email),
		)
		return &types.LoginResponse{}, ErrInvalidCredentials
	}

	s.logger.Info("user logged in successfully",
		zap.String("email", req.Email),
		zap.Uint32("id", u.ID),
	)

	// Gets existing session or a new one if one does not exist
	var store = s.cookieStore
	userSession, err := store.Get(ctx.Request, "user-session")
	if err != nil {
		s.logger.Error("failed to get session", zap.Error(err))
		return &types.LoginResponse{}, err
	}

	return &types.LoginResponse{
		Success:   true,
		Message:   "Login successful",
		User:      u,
		SessionID: userSession.ID,
	}, nil
}
