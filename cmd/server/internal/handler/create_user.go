package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"fakedating/pkg/model"
	"fakedating/pkg/util"
	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) CreateUser(w http.ResponseWriter, _ *http.Request) {
	user := generateUser()
	password := generateRandomPassword()

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		util.WriteErrorResponse("Failed to hash password", hashErr, http.StatusInternalServerError, w)
		return
	}
	user.PasswordHash = string(hashedPassword)

	persistedUser, createErr := h.userRepository.Create(user)
	if createErr != nil {
		util.WriteErrorResponse("Failed to persist new user", createErr, http.StatusInternalServerError, w)
		return
	}
	log.Printf("Created user %q", persistedUser.ID)

	encodedUser, marshalErr := json.Marshal(
		userWithOneTimePasswordReveal{
			User:     persistedUser,
			Password: password,
		},
	)
	if marshalErr != nil {
		util.WriteErrorResponse("Failed to marshal user to JSON", marshalErr, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(encodedUser)
}

func generateUser() model.User {
	gender := model.GenderFemale
	name := faker.FirstNameFemale()

	if rand.Float32() < 0.5 {
		gender = model.GenderMale
		name = faker.FirstNameMale()
	}

	name = fmt.Sprintf("%s %s", name, faker.LastName())

	return model.User{
		Email:  faker.Email(),
		Name:   name,
		Gender: gender,
		Age:    uint(18 + rand.Intn(50)),
	}
}

func generateRandomPassword() string {
	chars := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@£$%^#&*()")
	length := 15
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

type userWithOneTimePasswordReveal struct {
	model.User
	Password string
}
