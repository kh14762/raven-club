package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"raven-club/internal/pkg/types"
	"time"

	_ "github.com/lib/pq"
)

const ( // TODO: create a Config Struct that reads from a yaml file or something
	host     = "localhost"
	port     = 5432
	user     = "kev"
	password = "gulfstream"
	dbname   = "raven_club_db"
)

type Repository interface {
	Hook(lc fx.Lifecycle)
	CreateUser(ctx *gin.Context, user types.User) error
	GetUserByID(ctx *gin.Context, id string) (*types.User, error)
	GetUserByUsername(ctx *gin.Context, username string) (*types.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (*types.User, error)
	UpdateUser(ctx *gin.Context, user types.User) error
	DeleteUser(ctx *gin.Context, id string) error
	ListUsers(ctx *gin.Context) ([]types.User, error)
}

type repository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(logger *zap.Logger) Repository {
	return &repository{
		logger: logger,
		db:     nil,
	}
}

func (r *repository) Hook(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			r.logger.Info("Connecting to User db...")
			go func() {
				dbConn, err := r.NewConnection()
				r.db = dbConn
				if err != nil {
					r.logger.Error("Failed to connect to User db", zap.Error(err))
				}
				r.logger.Info("Successfully connected to User db")
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			r.logger.Info("Closing User db...")
			return r.db.Close()
		},
	})
}

func (r *repository) NewConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d "+
		"user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		r.logger.Error("failed to open database", zap.Error(err))
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func (r *repository) CreateUser(ctx *gin.Context, user types.User) error {
	query := `
		INSERT INTO users.users (username, email, password)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password)
	return err
}

func (r *repository) GetUserByID(ctx *gin.Context, id string) (*types.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
		WHERE id = $1
	`
	var user types.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByUsername(ctx *gin.Context, username string) (*types.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
		WHERE username = $1
	`
	var user types.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx *gin.Context, email string) (*types.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
		WHERE email = $1
	`
	var user types.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(ctx *gin.Context, user types.User) error {
	query := `
		UPDATE users.users
		SET username = $1, email = $2, password = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r *repository) DeleteUser(ctx *gin.Context, id string) error {
	query := `
		DELETE FROM users.users
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) ListUsers(ctx *gin.Context) ([]types.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Error("failed to close rows", zap.Error(err))
		}
	}(rows)

	var users []types.User
	for rows.Next() {
		var user types.User
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
