package domain

import "time"

type UserResponse struct {
	ID        string    `json:"userId"`
	NIP       int64     `json:"nip"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserQueryParams struct {
	ID        string `query:"userId"`
	Limit     string `query:"limit"`
	Offset    string `query:"offset"`
	Name      string `query:"name"`
	NIP       string `query:"nip"`
	Role      string `query:"role"`
	CreatedAt string `query:"createdAt"`
}
