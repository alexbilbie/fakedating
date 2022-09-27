package model

import (
	"github.com/segmentio/ksuid"
)

type Gender uint

const (
	GenderUnknown Gender = iota
	GenderFemale
	GenderMale
)

type User struct {
	ID                ksuid.KSUID
	Email             string
	PasswordHash      string `json:"-"`
	Name              string
	Gender            Gender
	Age               uint
	LocationLatitude  float64 `json:",omitempty"`
	LocationLongitude float64 `json:",omitempty"`
}
