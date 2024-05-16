package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserNurse interface {
	Register(ctx context.Context, db *sqlx.DB, nurse domain.UserNurse) error
	GetUserNurseByNIP(ctx context.Context, db *sqlx.DB, nip int64) (*domain.UserNurse, error)
	Update(ctx context.Context, db *sqlx.DB, nurse domain.UpdateUserNurse) error
}

type userNurse struct{}

func NewUserNurse() UserNurse {
	return &userNurse{}
}

func (un *userNurse) Register(ctx context.Context, db *sqlx.DB, nurse domain.UserNurse) error {
	query := `INSERT INTO users (id, created_at, nip, name, id_card_img, role)
	VALUES (:id, :created_at, :nip, :name, :id_card_img, :role)`
	_, err := db.NamedExecContext(ctx, query, nurse)
	if err != nil {
		return err
	}

	return nil
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

func (un *userNurse) Update(ctx context.Context, db *sqlx.DB, nurse domain.UpdateUserNurse) error {
	query := `UPDATE users
	SET nip = $2,
	name = $3
	WHERE id = $1`
	_, err := db.ExecContext(ctx, query, nurse.ID, nurse.NIP, nurse.Name)
	if err != nil {
		return err
	}

	return nil
}
