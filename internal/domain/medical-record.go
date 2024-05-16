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
	IdentityNumber int    `json:"identityNumber" validate:"required,intlen=16"`
	Symptoms       string `json:"symptoms" validate:"required,min=1,max=2000"`
	Medications    string `json:"medications" validate:"required,min=1,max=2000"`
}

type MedicalPatientRecord struct {
	IdentityNumber int    `json:"identityNumber"`
	PhoneNumber    string `json:"phoneNumber"`
	Name           string `json:"name"`
	BirthDate      string `json:"birthDate"`
	Gender         string `json:"gender"`
	IDCardImg      string `json:"identityCardScanImg"`
}

type UserRecord struct {
	ID   string `json:"userId"`
	NIP  string `json:"nip"`
	Name string `json:"name"`
}

type GetMedicalRecord struct {
	CreatedAt      time.Time            `json:"createdAt"`
	Symptoms       string               `json:"symptoms"`
	Medications    string               `json:"medications"`
	IdentityDetail MedicalPatientRecord `json:"identityDetail"`
	CreatedBy      UserRecord           `json:"createdBy"`
}

// Query Params
type MedicalPatientRecordQueryParams struct {
	IdentityNumber string `query:"identityNumber"`
}

type UserRecordQueryParams struct {
	ID  string `query:"userId"`
	NIP string `query:"nip"`
}

type MedicalRecordQueryParams struct {
	CreatedAt      string                          `query:"createdAt"`
	Limit          string                          `query:"limit"`
	Offset         string                          `query:"offset"`
	IdentityDetail MedicalPatientRecordQueryParams `query:"identityDetail"`
	CreatedBy      UserRecordQueryParams           `query:"createdBy"`
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
