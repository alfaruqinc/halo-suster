package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type MedicalPatient interface {
	Create(ctx context.Context, db *sqlx.DB, patient domain.MedicalPatient) error
}

type medicalPatient struct{}

func NewMedicalPatient() MedicalPatient {
	return &medicalPatient{}
}

func (mp *medicalPatient) Create(ctx context.Context, db *sqlx.DB, patient domain.MedicalPatient) error {
	query := `INSERT INTO medical_patients (id, created_at, identity_number, phone_number, name, birth_date, gender, id_card_img)
	VALUES (:id, :created_at, :identity_number, :phone_number, :name, :birth_date, :gender, :id_card_img)`
	_, err := db.NamedExecContext(ctx, query, patient)
	if err != nil {
		return err
	}

	return nil
}
