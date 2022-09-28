package payload

import (
	"fakedating/pkg/model"
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	User  model.User
	Token string
}
