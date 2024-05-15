package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserIT interface {
	Register(ctx context.Context, db *sqlx.DB, user domain.UserIT) error
	GetUserITByNIP(ctx context.Context, db *sqlx.DB, nip int64) (*domain.UserIT, error)
}

type userIT struct{}

func NewUserIT() UserIT {
	return &userIT{}
}

func (uit *userIT) Register(ctx context.Context, db *sqlx.DB, user domain.UserIT) error {
	query := `INSERT INTO users (id, created_at, nip, name, password, role)
	VALUES (:id, :created_at, :nip, :name, :password, :role)`
	_, err := db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}

func (uit *userIT) GetUserITByNIP(ctx context.Context, db *sqlx.DB, nip int64) (*domain.UserIT, error) {
	query := `SELECT id, created_at, nip, name, password, role
	FROM users
	WHERE nip = $1`
	var userIT domain.UserIT
	err := db.QueryRowxContext(ctx, query, nip).StructScan(&userIT)
	if err != nil {
		return nil, err
	}

	return &userIT, nil
}
