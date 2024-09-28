package model

type RegisterHTTPPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginHTTPPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
type LoginResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
}
