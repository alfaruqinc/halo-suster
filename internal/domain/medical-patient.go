package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type MedicalPatient struct {
	ID                  string    `db:"id"`
	CreatedAt           time.Time `db:"created_at"`
	IdentityNumber      int       `db:"identity_number"`
	PhoneNumber         string    `db:"phone_number"`
	Name                string    `db:"name"`
	BirthDate           string    `db:"birth_date"`
	Gender              string    `db:"gender"`
	IdentityCardScanImg string    `db:"id_card_img"`
}

type CreateMedicalPatient struct {
	IdentityNumber      int    `json:"identity_number"`
	PhoneNumber         string `json:"phone_number"`
	Name                string `json:"name"`
	BirthDate           string `json:"birth_date"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"id_card_img"`
}

type GetMedicalPatient struct {
	CreatedAt      time.Time `json:"createdAt"`
	IdentityNumber int       `json:"identity_number"`
	PhoneNumber    string    `json:"phone_number"`
	Name           string    `json:"name"`
	BirthDate      string    `json:"birth_date"`
	Gender         string    `json:"gender"`
}

func (cmp *CreateMedicalPatient) NewMedicalPatienFromDTO() MedicalPatient {
	id := ulid.Make()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	return MedicalPatient{
		ID:                  id.String(),
		CreatedAt:           createdAt,
		IdentityNumber:      cmp.IdentityNumber,
		PhoneNumber:         cmp.PhoneNumber,
		Name:                cmp.Name,
		BirthDate:           cmp.BirthDate,
		Gender:              cmp.Gender,
		IdentityCardScanImg: cmp.IdentityCardScanImg,
	}
}