package session

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const ( // TODO: create a Config Struct that reads from a yaml file or something
	host     = "localhost"
	port     = 5432
	user     = "kev"
	password = "gulfstream"
	dbname   = "session_db"
)

type Repository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(logger *zap.Logger) *Repository {
	return &Repository{
		logger: logger,
		db:     nil,
	}
}

func (r *Repository) Hook(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			r.logger.Info("Connecting to test db...")
			go func() {
				dbConn, err := r.NewConnection()
				r.db = dbConn
				if err != nil {
					r.logger.Error("Failed to connect to test db", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			r.logger.Info("Closing test db...")
			return r.db.Close()
		},
	})
}

func (r *Repository) NewConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d "+
		"user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		r.logger.Error("failed to connect to database", zap.Error(err))
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	r.logger.Info("Successfully connected to test db")
	return db, nil
}

//func (r *repository) GetDb() (*sql.DB, error) {
//	if r.db == nil {
//		r.logger.Info("Database connection is empty, initializing database...")
//		r.db, _ = r.NewConnection()
//	}
//	return r.db, nil
//}
