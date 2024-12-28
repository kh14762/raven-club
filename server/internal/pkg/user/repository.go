package user

import (
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"raven-club/internal/pkg/types"

	_ "github.com/lib/pq"
)

const ( // TODO: create a Config Struct that reads from a yaml file or something
	host     = "localhost"
	port     = 5432
	user     = "kev"
	password = "gulfstream"
	dbname   = "raven_club_db"
)

// Repository TODO: define crud ops
type Repository interface {
	CreateUser(ctx context.Context, user types.User) error
	GetUserByID(ctx context.Context, id string) (*types.User, error)
	GetUserByUsername(ctx context.Context, username string) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	UpdateUser(ctx context.Context, user types.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]types.User, error)
}

type repository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(logger *zap.Logger, db *sql.DB) Repository {
	return &repository{
		logger: logger,
		db:     db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user types.User) error {
	query := `
		INSERT INTO users.users (username, email, password)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password)
	return err
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*types.User, error) {
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

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
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

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
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

func (r *repository) UpdateUser(ctx context.Context, user types.User) error {
	query := `
		UPDATE users.users
		SET username = $1, email = $2, password = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r *repository) DeleteUser(ctx context.Context, id string) error {
	query := `
		DELETE FROM users.users
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) ListUsers(ctx context.Context) ([]types.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users.users
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
