package list_projects

import (
	"context"

	"github.com/google/uuid"

	"app/core/project/domain"
)

type Query struct{ OwnerID uuid.UUID }

type Handler struct{ repo domain.Repository }

func NewHandler(r domain.Repository) *Handler { return &Handler{repo: r} }

func (h *Handler) Handle(ctx context.Context, q Query) ([]*domain.Project, error) {
	return h.repo.List(ctx, q.OwnerID)
}
