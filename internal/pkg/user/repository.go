package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

const ( // TODO: create a Config Struct that reads from a yaml file or something
	host     = "localhost"
	port     = 5432
	user     = "kev"
	password = "gulfstream"
	dbname   = "raven_club_db"
)

type Repository interface {
	NewConnection(lc fx.Lifecycle)
	CreateUser(ctx *gin.Context, user User) error
	GetUserByID(ctx *gin.Context, id string) (*User, error)
	GetUserByUsername(ctx *gin.Context, username string) (*User, error)
	GetUserByEmail(ctx *gin.Context, email string) (*User, error)
	UpdateUser(ctx *gin.Context, user User) error
	DeleteUser(ctx *gin.Context, id string) error
	ListUsers(ctx *gin.Context) ([]User, error)
}

type repository struct {
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func NewRepository(logger *zap.Logger) Repository {
	return &repository{
		logger: logger,
		pool:   nil,
	}
}

func (r *repository) NewConnection(lc fx.Lifecycle) {
	dbURL := os.Getenv("USER_DATABASE_URL")
	if dbURL == "" {
		err := errors.New("USER_DATABASE_URL environment variable not set")
		r.logger.Error("", zap.Error(err))
	}

	// Create a context with a timeout for connecting to the database.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a connection pool. Adjust configuration as required.
	//	db.SetMaxIdleConns(10)
	//	db.SetMaxOpenConns(100)
	//	db.SetConnMaxLifetime(time.Hour) // TODO: set these parameters
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		r.logger.Error("unable to create connection pool: ", zap.Error(err))
	}

	r.pool = pool

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			r.logger.Info("Starting user database connection pool...")
			var currentTime time.Time
			err := r.pool.QueryRow(ctx, "SELECT NOW()").Scan(&currentTime)
			if err != nil {
				return fmt.Errorf("failed to verify connection: %w", err)
			}
			log.Printf("Database connection pool started. Current DB time: %v", currentTime)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			r.logger.Info("Closing user database connection pool...")
			r.pool.Close()
			return nil
		},
	})
}

func (r *repository) CreateUser(ctx *gin.Context, user User) error {
	query := `
		INSERT INTO users.users (username, email, password)
		VALUES ($1, $2, $3)
	`
	_, err := r.pool.Exec(ctx, query, user.Username, user.Email, user.Password)
	return err
}

func (r *repository) GetUserByID(ctx *gin.Context, id string) (*User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
		WHERE id = $1
	`
	var user User
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByUsername(ctx *gin.Context, username string) (*User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
		WHERE username = $1
	`
	var user User
	err := r.pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx *gin.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
		WHERE email = $1
	`
	var user User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(ctx *gin.Context, user User) error {
	query := `
		UPDATE users.users
		SET username = $1, email = $2, password = $3
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r *repository) DeleteUser(ctx *gin.Context, id string) error {
	query := `
		DELETE FROM users.users
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *repository) ListUsers(ctx *gin.Context) ([]User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // ensure rows are closed

	var users []User
	for rows.Next() {
		var user User
		// Scan the row into the user struct fields
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
