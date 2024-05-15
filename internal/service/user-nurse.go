package service

import (
	"context"
	"database/sql"
	"health-record/internal/domain"
	"health-record/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserNurse interface {
	Login(ctx context.Context, user domain.LoginUserNurse) (*domain.UserNurseResponse, domain.ErrorMessage)
}

type userNurse struct {
	db            *sqlx.DB
	jwtSecret     string
	bcryptSalt    string
	userNurseRepo repository.UserNurse
}

func NewUserNurse(db *sqlx.DB, jwtSecret string, bcryptSalt string, userNurseRepo repository.UserNurse) UserNurse {
	return &userNurse{
		db:            db,
		jwtSecret:     jwtSecret,
		bcryptSalt:    bcryptSalt,
		userNurseRepo: userNurseRepo,
	}
}

func (un *userNurse) Login(ctx context.Context, user domain.LoginUserNurse) (*domain.UserNurseResponse, domain.ErrorMessage) {
	isNurseRole := (user.NIP / 1e10) == 303
	if !isNurseRole {
		return nil, domain.NewErrNotFound("user is not found")
	}

	userNurse, err := un.userNurseRepo.GetUserNurseByNIP(ctx, un.db, user.NIP)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewErrNotFound("user is not found")
		}
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	if !userNurse.Password.Valid {
		return nil, domain.NewErrBadRequest("user does not have access")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userNurse.Password.String), []byte(user.Password))
	if err != nil {
		return nil, domain.NewErrNotFound("password is wrong")
	}

	claims := jwt.MapClaims{
		"id":   userNurse.ID,
		"nip":  userNurse.NIP,
		"role": userNurse.Role,
		"exp":  time.Now().Add(time.Hour * 8).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(un.jwtSecret))
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	loginResp := domain.UserNurseResponse{
		ID:          userNurse.ID,
		NIP:         userNurse.NIP,
		Name:        userNurse.Name,
		AccessToken: t,
	}

	return &loginResp, nil
}