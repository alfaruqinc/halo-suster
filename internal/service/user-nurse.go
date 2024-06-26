package service

import (
	"context"
	"database/sql"
	"health-record/internal/domain"
	"health-record/internal/helper"
	"health-record/internal/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserNurse interface {
	Register(ctx context.Context, nurse domain.UserNurse) (*domain.UserNurseResponse, domain.ErrorMessage)
	Login(ctx context.Context, user domain.LoginUserNurse) (*domain.UserNurseResponse, domain.ErrorMessage)
	Update(ctx context.Context, nurse domain.UpdateUserNurse) domain.ErrorMessage
	Delete(ctx context.Context, nurseId string) domain.ErrorMessage
	GiveAccess(ctx context.Context, nurse domain.AccessSystemUserNurse) domain.ErrorMessage
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

func (un *userNurse) Register(ctx context.Context, nurse domain.UserNurse) (*domain.UserNurseResponse, domain.ErrorMessage) {
	err := un.userNurseRepo.Register(ctx, un.db, nurse)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				return nil, domain.NewErrConflict("nip already exists")
			}
		}
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	registerResp := domain.UserNurseResponse{
		ID:   nurse.ID,
		NIP:  nurse.NIP,
		Name: nurse.Name,
	}

	return &registerResp, nil
}

func (un *userNurse) Login(ctx context.Context, user domain.LoginUserNurse) (*domain.UserNurseResponse, domain.ErrorMessage) {
	isNurseRole := helper.MatchRole(user.NIP, domain.NurseRole)
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

func (un *userNurse) Update(ctx context.Context, nurse domain.UpdateUserNurse) domain.ErrorMessage {
	// isNurseRole := helper.MatchRole(nurse.NIP, domain.NurseRole)
	// if !isNurseRole {
	// 	return domain.NewErrNotFound("user is not nurse")
	// }

	affRow, err := un.userNurseRepo.Update(ctx, un.db, nurse)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				return domain.NewErrConflict("nip already exists")
			}
		}
		return domain.NewErrInternalServerError(err.Error())
	}
	if affRow == 0 {
		return domain.NewErrNotFound("nurse is not found")
	}

	return nil
}

func (un *userNurse) Delete(ctx context.Context, nurseId string) domain.ErrorMessage {
	affRow, err := un.userNurseRepo.Delete(ctx, un.db, nurseId)
	if err != nil {
		return domain.NewErrInternalServerError(err.Error())
	}
	if affRow == 0 {
		return domain.NewErrNotFound("nurse is not found")
	}

	return nil
}

func (un *userNurse) GiveAccess(ctx context.Context, nurse domain.AccessSystemUserNurse) domain.ErrorMessage {
	salt, err := strconv.Atoi(un.bcryptSalt)
	if err != nil {
		return domain.NewErrInternalServerError(err.Error())
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(nurse.Password), salt)
	if err != nil {
		return domain.NewErrInternalServerError(err.Error())
	}
	nurse.Password = string(hashedPassword)

	affRow, err := un.userNurseRepo.GiveAccess(ctx, un.db, nurse)
	if err != nil {
		return domain.NewErrInternalServerError(err.Error())
	}
	if affRow == 0 {
		return domain.NewErrNotFound("nurse is not found")
	}

	return nil
}
