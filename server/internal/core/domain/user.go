package domain

import "time"

type User struct {
	ID           string    `json:"id,omitempty"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}
