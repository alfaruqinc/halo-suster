package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type MedicalPatient struct {
	ID             string    `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	IdentityNumber int       `db:"identity_number"`
	PhoneNumber    string    `db:"phone_number"`
	Name           string    `db:"name"`
	BirthDate      string    `db:"birth_date"`
	Gender         string    `db:"gender"`
	IDCardImg      string    `db:"id_card_img"`
}

type CreateMedicalPatient struct {
	IdentityNumber int    `json:"identityNumber" validate:"required,intlen=16"`
	PhoneNumber    string `json:"phoneNumber" validate:"required,startswith=+62,min=10,max=15"`
	Name           string `json:"name" validate:"required,min=3,max=30"`
	BirthDate      string `json:"birthDate" validate:"required,iso8601"`
	Gender         string `json:"gender" validate:"required,oneof=male female"`
	IDCardImg      string `json:"identityCardScanImg" validate:"required,url"`
}

type GetMedicalPatient struct {
	CreatedAt      time.Time `json:"createdAt"`
	IdentityNumber int       `json:"identityNumber"`
	PhoneNumber    string    `json:"phoneNumber"`
	Name           string    `json:"name"`
	BirthDate      string    `json:"birthDate"`
	Gender         string    `json:"gender"`
}

type MedicalPatientQueryParams struct {
	IdentityNumber string `query:"identityNumber"`
	Limit          string `query:"limit"`
	Offset         string `query:"offset"`
	Name           string `query:"name"`
	PhoneNumber    string `query:"phoneNumber"`
	CreatedAt      string `query:"createdAt"`
}

func (cmp *CreateMedicalPatient) NewMedicalPatienFromDTO() MedicalPatient {
	id := ulid.Make()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	return MedicalPatient{
		ID:             id.String(),
		CreatedAt:      createdAt,
		IdentityNumber: cmp.IdentityNumber,
		PhoneNumber:    cmp.PhoneNumber,
		Name:           cmp.Name,
		BirthDate:      cmp.BirthDate,
		Gender:         cmp.Gender,
		IDCardImg:      cmp.IDCardImg,
	}
}
