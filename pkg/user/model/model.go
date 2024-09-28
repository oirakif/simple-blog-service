package model

import "time"

type FindUserFilterQuery struct {
	ID           *int
	Name         *string
	Email        *string
	PasswordHash *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type User struct {
	ID           int       `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"passwordHash,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}
