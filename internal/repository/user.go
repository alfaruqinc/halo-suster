package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type User interface {
	GetAllUser(ctx context.Context, db *sqlx.DB) ([]domain.UserResponse, error)
}

type user struct{}

func NewUser() User {
	return &user{}
}

func (u *user) GetAllUser(ctx context.Context, db *sqlx.DB) ([]domain.UserResponse, error) {
	query := `SELECT id, created_at, nip, name
	FROM users`
	users := []domain.UserResponse{}
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u domain.UserResponse
		err = rows.Scan(&u.ID, &u.CreatedAt, &u.NIP, &u.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}
