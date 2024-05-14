package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type MedicalRecord struct {
	ID               string    `db:"id"`
	CreatedAt        time.Time `db:"created_at"`
	IdentityNumber   int       `db:"identity_number"`
	Symptoms         string    `db:"symptoms"`
	Medications      string    `db:"medications"`
	MedicalPatientID string    `db:"medical_patient_id"`
	CreatedByID      string    `db:"created_by_id"`
}

type CreateMedicalRecord struct {
	IdentityNumber int    `json:"identityNumber"`
	Symptoms       string `json:"symptoms"`
	Medications    string `json:"medications"`
}

type GetMedicalRecord struct {
	CreatedAt      time.Time      `json:"createdAt"`
	Symptoms       string         `json:"symptoms"`
	Medications    string         `json:"medications"`
	IdentityDetail MedicalPatient `json:"identityDetail"`
	CreatedBy      UserNurse      `json:"createdBy"`
}

func (cmr *CreateMedicalRecord) NewMedicalRecordFromDTO() MedicalRecord {
	id := ulid.Make()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	return MedicalRecord{
		ID:             id.String(),
		CreatedAt:      createdAt,
		IdentityNumber: cmr.IdentityNumber,
		Symptoms:       cmr.Symptoms,
		Medications:    cmr.Medications,
	}
}
