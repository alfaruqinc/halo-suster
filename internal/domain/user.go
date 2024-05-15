package domain

import "time"

type UserResponse struct {
	ID        string    `json:"userId"`
	NIP       int64     `json:"nip"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
