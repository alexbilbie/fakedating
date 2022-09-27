package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fakedating/pkg/model"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateUser(t *testing.T) {
	h := Handler{userRepository: mockCreateUserRepository{}}

	w := &httptest.ResponseRecorder{}
	h.CreateUser(w, nil)

	assert.Equal(t, w.Result().StatusCode, http.StatusCreated)
}

type mockCreateUserRepository struct {
	UserRepository
}

func (mockCreateUserRepository) Create(u model.User) (model.User, error) {
	u.ID = ksuid.New()
	return u, nil
}
