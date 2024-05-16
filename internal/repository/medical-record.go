package repository

import (
	"context"
	"health-record/internal/domain"

	"github.com/jmoiron/sqlx"
)

type MedicalRecord interface {
	Create(ctx context.Context, db *sqlx.DB, record domain.MedicalRecord) (int64, error)
	GetAllMedicalRecords(ctx context.Context, db *sqlx.DB) ([]domain.GetMedicalRecord, error)
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

func (mr *medicalRecord) GetAllMedicalRecords(ctx context.Context, db *sqlx.DB) ([]domain.GetMedicalRecord, error) {
	query := `SELECT mr.created_at, mr.symptoms, mr.medications,
	mp.identity_number, mp.phone_number, mp.name, mp.birth_date, mp.gender, mp.id_card_img,
	u.id, u.nip, u.name
	FROM medical_records mr
	JOIN medical_patients mp ON mp.id = mr.medical_patient_id
	JOIN users u ON u.id = mr.created_by_id`
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	records := []domain.GetMedicalRecord{}
	for rows.Next() {
		r := domain.GetMedicalRecord{}
		err := rows.Scan(
			&r.CreatedAt, &r.Symptoms, &r.Medications,
			&r.IdentityDetail.IdentityNumber, &r.IdentityDetail.PhoneNumber, &r.IdentityDetail.Name,
			&r.IdentityDetail.BirthDate, &r.IdentityDetail.Gender, &r.IdentityDetail.IDCardImg,
			&r.CreatedBy.ID, &r.CreatedBy.NIP, &r.CreatedBy.Name,
		)
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}
