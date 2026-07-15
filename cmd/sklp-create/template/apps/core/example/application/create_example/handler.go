// Package create_example is the write-side use case: create a new Example.
package create_example

import (
	"context"

	"github.com/google/uuid"

	"app/core/example/domain"
	"app/core/example/infrastructure"
)

// Command is the create-example command.
type Command struct {
	Name string
}

// Result carries the new aggregate id.
type Result struct {
	ID domain.ID
}

// Handler executes the command against the aggregate + repository.
type Handler struct {
	repo *infrastructure.Repository
}

func NewHandler(repo *infrastructure.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	id := uuid.NewString()
	e := domain.New(id)
	if err := e.Create(cmd.Name); err != nil {
		return Result{}, err
	}
	if err := h.repo.Save(ctx, e); err != nil {
		return Result{}, err
	}
	return Result{ID: id}, nil
}
