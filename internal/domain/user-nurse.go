package domain

import (
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"
)

type UserNurse struct {
	ID        string         `db:"id"`
	CreatedAt time.Time      `db:"created_at"`
	NIP       int64          `db:"nip"`
	Name      string         `db:"name"`
	IDCardImg string         `db:"id_card_img"`
	Password  sql.NullString `db:"password"`
	Role      string         `db:"role"`
}

type RegisterUserNurse struct {
	NIP       int64  `json:"nip" validate:"required,nip=nurse"`
	Name      string `json:"name" validate:"required,min=5,max=50"`
	IDCardImg string `json:"identityCardScanImg" validate:"required,url"`
}

type AccessSystemUserNurse struct {
	ID       string `json:"id"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type LoginUserNurse struct {
	NIP      int64  `json:"nip" validate:"required,nip=nurse"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type UserNurseResponse struct {
	ID          string `json:"userId"`
	NIP         int64  `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken,omitempty"`
}

type UpdateUserNurse struct {
	ID   string `json:"id"`
	NIP  int64  `json:"nip" validate:"required,nip=nurse"`
	Name string `json:"name" validate:"required,min=5,max=50"`
}

type GetUserNurse struct {
	ID        string    `json:"userId"`
	NIP       int       `json:"nip"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u RegisterUserNurse) NewUserNurseFromDTO() UserNurse {
	id := ulid.Make()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	return UserNurse{
		ID:        id.String(),
		CreatedAt: createdAt,
		NIP:       u.NIP,
		Name:      u.Name,
		IDCardImg: u.IDCardImg,
		Role:      "nurse",
	}
}
