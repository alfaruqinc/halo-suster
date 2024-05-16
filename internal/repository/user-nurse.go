package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserNurse interface {
	Register(ctx context.Context, db *sqlx.DB, nurse domain.UserNurse) error
	GetUserNurseByNIP(ctx context.Context, db *sqlx.DB, nip int64) (*domain.UserNurse, error)
	Update(ctx context.Context, db *sqlx.DB, nurse domain.UpdateUserNurse) (int64, error)
	Delete(ctx context.Context, db *sqlx.DB, nurseId string) (int64, error)
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

func (un *userNurse) Update(ctx context.Context, db *sqlx.DB, nurse domain.UpdateUserNurse) (int64, error) {
	query := `UPDATE users
	SET nip = $2,
	name = $3
	WHERE id = $1`
	res, err := db.ExecContext(ctx, query, nurse.ID, nurse.NIP, nurse.Name)
	if err != nil {
		return 0, err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return aff, nil
}

func (un *userNurse) Delete(ctx context.Context, db *sqlx.DB, nurseId string) (int64, error) {
	query := `DELETE FROM users
	WHERE id = $1
	AND CAST((nip / 1e10) AS VARCHAR(3)) = '303'`
	res, err := db.ExecContext(ctx, query, nurseId)
	if err != nil {
		return 0, err
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affRow, nil
}
