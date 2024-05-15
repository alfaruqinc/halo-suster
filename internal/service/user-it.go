package service

import (
	"context"
	"health-record/internal/domain"
	"health-record/internal/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserIT interface {
	Register(ctx context.Context, user domain.UserIT) (*domain.UserITResponse, domain.ErrorMessage)
}

type userIT struct {
	db             *sqlx.DB
	userRepository repository.UserIT
	jwtSecret      string
	bcryptSalt     string
}

func NewUserIT(db *sqlx.DB, jwtSecret string, bcryptSalt string, userRepository repository.UserIT) UserIT {
	return &userIT{
		db:             db,
		jwtSecret:      jwtSecret,
		bcryptSalt:     bcryptSalt,
		userRepository: userRepository,
	}
}

func (uit *userIT) Register(ctx context.Context, user domain.UserIT) (*domain.UserITResponse, domain.ErrorMessage) {
	salt, err := strconv.Atoi(uit.bcryptSalt)
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), salt)
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}
	user.Password = string(hashedPassword)

	err = uit.userRepository.Register(ctx, uit.db, user)
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	claims := jwt.MapClaims{
		"id":   user.ID,
		"nip":  user.NIP,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 8).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(uit.jwtSecret))
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	registerResp := domain.UserITResponse{
		ID:          user.ID,
		NIP:         user.NIP,
		Name:        user.Name,
		AccessToken: t,
	}

	return &registerResp, nil
}
