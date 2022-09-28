package payload

import (
	"fakedating/pkg/model"
)

type CreateUserResponse struct {
	model.User
	Password string
}
