package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type UserNurse struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	NIP       int       `db:"nip"`
	Name      string    `db:"name"`
	IDCardImg string    `db:"id_card_img"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
}

type RegisterUserNurse struct {
	NIP       int    `json:"nip"`
	Name      string `json:"name"`
	IDCardImg string `json:"identityCardScanImg"`
}

type AccessSystemUserNurse struct {
	Password string `json:"password"`
}

type LoginUserNurse struct {
	NIP      int    `json:"nip"`
	Password string `json:"password"`
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
