package core

import (
	"github.com/google/uuid"
)

type User struct {
	UUID     uuid.UUID `json:"UUID"`
	Username string    `json:"username"`
	Slugs    []Slug    `json:"slugs"`
}

type UserPut struct {
	Add_slugs    []string `json:"add_slugs"`
	Delete_slugs []string `json:"delete_slugs"`
}

type UserRequestCreate struct {
	Username string `json:"username"`
}
