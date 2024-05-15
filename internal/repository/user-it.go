package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserIT interface {
	Register(ctx context.Context, db *sqlx.DB, user domain.UserIT) error
}

type userIT struct{}

func NewUserIT() UserIT {
	return &userIT{}
}

func (u *userIT) Register(ctx context.Context, db *sqlx.DB, user domain.UserIT) error {
	query := `INSERT INTO users (id, created_at, nip, name, password, role)
	VALUES (:id, :created_at, :nip, :name, :password, :role)`
	_, err := db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}
