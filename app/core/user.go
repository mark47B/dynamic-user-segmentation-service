package core

import (
	"github.com/google/uuid"
)

type User struct {
	UUID     uuid.UUID `json:"UUID"`
	Username string    `json:"username"`
	Slugs    []Slug    `json:"slugs"`
}
