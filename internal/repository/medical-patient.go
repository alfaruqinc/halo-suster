package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type MedicalPatient interface {
	Create(ctx context.Context, db *sqlx.DB, patient domain.MedicalPatient) error
	GetAllMedicalPatients(ctx context.Context, db *sqlx.DB) ([]domain.GetMedicalPatient, error)
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

func (mp *medicalPatient) GetAllMedicalPatients(ctx context.Context, db *sqlx.DB) ([]domain.GetMedicalPatient, error) {
	query := `SELECT created_at, identity_number, phone_number, name, birth_date, gender
	FROM medical_patients`
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	medicalPatients := []domain.GetMedicalPatient{}
	for rows.Next() {
		var mp domain.GetMedicalPatient
		err := rows.Scan(
			&mp.CreatedAt, &mp.IdentityNumber, &mp.PhoneNumber,
			&mp.Name, &mp.BirthDate, &mp.Gender,
		)
		if err != nil {
			return nil, err
		}
		medicalPatients = append(medicalPatients, mp)
	}

	return medicalPatients, nil
}
