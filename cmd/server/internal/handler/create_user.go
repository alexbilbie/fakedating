package handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"fakedating/pkg/model"
	"fakedating/pkg/payload"
	"fakedating/pkg/util"
	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) CreateUser(w http.ResponseWriter, _ *http.Request) {
	user := generateUser()
	password := generateRandomPassword()

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		log.Printf("Failed to hash password: %v", hashErr)
		util.WriteErrorResponse("Failed to hash password", hashErr, http.StatusInternalServerError, w)
		return
	}
	user.PasswordHash = string(hashedPassword)

	persistedUser, createErr := h.userRepository.Create(user)
	if createErr != nil {
		log.Printf("Failed to persist new user: %v", createErr)
		util.WriteErrorResponse("Failed to persist new user", createErr, http.StatusInternalServerError, w)
		return
	}
	log.Printf("Created user %q", persistedUser.ID)

	util.WriteJSONResponse(
		payload.CreateUserResponse{
			User:     persistedUser,
			Password: password,
		},
		http.StatusCreated,
		w,
	)
}

func generateUser() model.User {
	gender := model.GenderFemale
	name := faker.FirstNameFemale()

	if rand.Float32() < 0.5 {
		gender = model.GenderMale
		name = faker.FirstNameMale()
	}

	name = fmt.Sprintf("%s %s", name, faker.LastName())

	// Generate a fake location within London
	latMin := 51.416639
	latMax := 51.627694
	lat := latMin + rand.Float64()*(latMax-latMin)

	longMin := -0.367440
	longMax := 0.062400
	long := longMin + rand.Float64()*(longMax-longMin)

	return model.User{
		Email:    strings.ToLower(faker.Email()),
		Name:     name,
		Gender:   gender,
		Age:      uint(18 + rand.Intn(50)),
		Location: model.Location{Latitude: lat, Longitude: long},
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
