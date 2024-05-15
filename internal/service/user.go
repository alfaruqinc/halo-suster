package service

import (
	"context"
	"health-record/internal/domain"
	"health-record/internal/repository"

	"github.com/jmoiron/sqlx"
)

type User interface {
	GetAllUsers(ctx context.Context) ([]domain.UserResponse, domain.ErrorMessage)
}

type user struct {
	db       *sqlx.DB
	userRepo repository.User
}

func NewUser(db *sqlx.DB, userRepo repository.User) User {
	return &user{
		db:       db,
		userRepo: userRepo,
	}
}

func (u *user) GetAllUsers(ctx context.Context) ([]domain.UserResponse, domain.ErrorMessage) {
	users, err := u.userRepo.GetAllUser(ctx, u.db)
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	return users, nil
}
