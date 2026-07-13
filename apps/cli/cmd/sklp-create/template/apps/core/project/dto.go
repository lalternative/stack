package project

import (
	"time"

	"github.com/google/uuid"

	"app/core/project/domain"
)

type CreateRequest struct {
	Name string `json:"name"`
}

type ProjectDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

func toDTO(p *domain.Project) ProjectDTO {
	return ProjectDTO{ID: p.ID, Name: p.Name, OwnerID: p.OwnerID, CreatedAt: p.CreatedAt}
}
