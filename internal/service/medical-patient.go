package service

import (
	"context"
	"health-record/internal/domain"
	"health-record/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type MedicalPatient interface {
	Create(ctx context.Context, patient domain.MedicalPatient) domain.ErrorMessage
	GetAllMedicalPatients(ctx context.Context) ([]domain.GetMedicalPatient, domain.ErrorMessage)
}

type medicalPatient struct {
	db                 *sqlx.DB
	medicalPatientRepo repository.MedicalPatient
}

func NewMedicalPatient(db *sqlx.DB, medicalPatientRepo repository.MedicalPatient) MedicalPatient {
	return &medicalPatient{
		db:                 db,
		medicalPatientRepo: medicalPatientRepo,
	}
}

func (mp *medicalPatient) Create(ctx context.Context, patient domain.MedicalPatient) domain.ErrorMessage {
	err := mp.medicalPatientRepo.Create(ctx, mp.db, patient)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				return domain.NewErrConflict("identity number already exists")
			}
		}
		return domain.NewErrInternalServerError(err.Error())
	}

	return nil
}

func (mp *medicalPatient) GetAllMedicalPatients(ctx context.Context) ([]domain.GetMedicalPatient, domain.ErrorMessage) {
	patients, err := mp.medicalPatientRepo.GetAllMedicalPatients(ctx, mp.db)
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	return patients, nil
}
