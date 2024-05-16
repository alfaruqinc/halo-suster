package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type MedicalRecord interface {
	Create(ctx context.Context, db *sqlx.DB, record domain.MedicalRecord) (int64, error)
}

type medicalRecord struct{}

func NewMedicalRecord() MedicalRecord {
	return &medicalRecord{}
}

func (mr *medicalRecord) Create(ctx context.Context, db *sqlx.DB, record domain.MedicalRecord) (int64, error) {
	query := `INSERT INTO medical_records (medical_patient_id, id, created_at, identity_number, symptoms, medications, created_by_id)
	SELECT mp.id, :id, :created_at, :identity_number, :symptoms, :medications, :created_by_id
	FROM medical_patients mp
	WHERE identity_number = :identity_number`
	res, err := db.NamedExecContext(ctx, query, record)
	if err != nil {
		return 0, err
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affRow, nil
}
