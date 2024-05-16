package service

import (
	"context"
	"health-record/internal/domain"
	"health-record/internal/repository"

	"github.com/jmoiron/sqlx"
)

type MedicalRecord interface {
	Create(ctx context.Context, record domain.MedicalRecord) domain.ErrorMessage
	GetAllMedicalRecords(ctx context.Context) ([]domain.GetMedicalRecord, domain.ErrorMessage)
}

type medicalRecord struct {
	db                *sqlx.DB
	medicalRecordRepo repository.MedicalRecord
}

func NewMedicalRecord(db *sqlx.DB, medicalRecordRepo repository.MedicalRecord) MedicalRecord {
	return &medicalRecord{
		db:                db,
		medicalRecordRepo: medicalRecordRepo,
	}
}

func (mr *medicalRecord) Create(ctx context.Context, record domain.MedicalRecord) domain.ErrorMessage {
	affRow, err := mr.medicalRecordRepo.Create(ctx, mr.db, record)
	if err != nil {
		return domain.NewErrInternalServerError(err.Error())
	}
	if affRow == 0 {
		return domain.NewErrBadRequest("identity number is not exists")
	}

	return nil
}

func (mr *medicalRecord) GetAllMedicalRecords(ctx context.Context) ([]domain.GetMedicalRecord, domain.ErrorMessage) {
	records, err := mr.medicalRecordRepo.GetAllMedicalRecords(ctx, mr.db)
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	return records, nil
}
