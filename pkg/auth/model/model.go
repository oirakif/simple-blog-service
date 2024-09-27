package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterHTTPPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginHTTPPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID           int       `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"passwordHash,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}

type RegisterResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Token   string      `json:"token,omitempty"`
}

type FindUserFilterQuery struct {
	ID           *int       `json:"id,omitempty"`
	Name         *string    `json:"name,omitempty"`
	Email        *string    `json:"email,omitempty"`
	PasswordHash *string    `json:"passwordHash,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
