package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("project not found")

type Project struct {
	ID        uuid.UUID
	Name      string
	OwnerID   uuid.UUID
	CreatedAt time.Time
}

func New(name string, owner uuid.UUID) (*Project, error) {
	if name == "" {
		return nil, errors.New("name required")
	}
	return &Project{
		ID:        uuid.New(),
		Name:      name,
		OwnerID:   owner,
		CreatedAt: time.Now().UTC(),
	}, nil
}
