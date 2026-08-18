// Package user
package user

import "github.com/google/uuid"

type User struct {
	Name string    `yaml:"name"`
	ID   uuid.UUID `yaml:"id"`
}
