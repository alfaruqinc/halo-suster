package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserNurse interface {
	GetUserNurseByNIP(ctx context.Context, db *sqlx.DB, nip int64) (*domain.UserNurse, error)
}

type userNurse struct{}

func NewUserNurse() UserNurse {
	return &userNurse{}
}

func (un *userNurse) GetUserNurseByNIP(ctx context.Context, db *sqlx.DB, nip int64) (*domain.UserNurse, error) {
	query := `SELECT id, created_at, nip, name, id_card_img, password, role
	FROM users
	WHERE nip = $1`
	var userNurse domain.UserNurse
	err := db.QueryRowxContext(ctx, query, nip).StructScan(&userNurse)
	if err != nil {
		return nil, err
	}

	return &userNurse, nil
}
