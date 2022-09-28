package payload

import (
	"fakedating/pkg/model"
)

type ListProfilesResponse struct {
	Matches []model.User `json:",omitempty"`
}
