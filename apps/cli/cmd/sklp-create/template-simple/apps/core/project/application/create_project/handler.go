package create_project

import (
	"context"

	"github.com/google/uuid"

	"app/core/project/domain"
)

type Command struct {
	Name    string
	OwnerID uuid.UUID
}

type Handler struct{ repo domain.Repository }

func NewHandler(r domain.Repository) *Handler { return &Handler{repo: r} }

func (h *Handler) Handle(ctx context.Context, cmd Command) (*domain.Project, error) {
	p, err := domain.New(cmd.Name, cmd.OwnerID)
	if err != nil {
		return nil, err
	}
	if err := h.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
