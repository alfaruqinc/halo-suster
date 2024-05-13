package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type UserIT struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Nip       int       `db:"nip"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
}

type RegisterUserIT struct {
	Nip      int    `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginUserIT struct {
	Nip      int    `json:"nip"`
	Password string `json:"password"`
}

type GetUserIT struct {
	ID        string    `json:"userId"`
	Nip       int       `json:"nip"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u RegisterUserIT) NewUserITFromDTO() UserIT {
	id := ulid.Make()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	return UserIT{
		ID:        id.String(),
		CreatedAt: createdAt,
		Nip:       u.Nip,
		Name:      u.Name,
		Password:  u.Password,
		Role:      "it",
	}
}
