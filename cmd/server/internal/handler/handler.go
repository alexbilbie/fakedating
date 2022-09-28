package handler

import (
	"net/http"

	"fakedating/pkg/model"
	"fakedating/pkg/util"
	"github.com/segmentio/ksuid"
)

type UserRepository interface {
	Create(model.User) (model.User, error)
	GetByID(ksuid.KSUID) (model.User, error)
	GetByEmail(email string) (model.User, error)
	ListMatches(userID ksuid.KSUID, options ...model.SearchParameterOpt) ([]model.User, error)
	SaveSwipe(swiperID ksuid.KSUID, recipient ksuid.KSUID, matched bool) (model.ProfileMatch, error)
}

type AuthRepository interface {
	CreateTokenForUser(ksuid.KSUID) (string, error)
}

func New(authRepo AuthRepository, userRepo UserRepository) Handler {
	return Handler{
		authRepository: authRepo,
		userRepository: userRepo,
	}
}

type Handler struct {
	authRepository AuthRepository
	userRepository UserRepository
}

type InvalidRoute struct{}

func (InvalidRoute) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	util.WriteErrorResponse("Unknown API route", nil, http.StatusNotFound, w)
}
