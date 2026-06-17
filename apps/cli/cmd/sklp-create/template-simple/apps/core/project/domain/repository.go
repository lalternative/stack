package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, p *Project) error
	List(ctx context.Context, owner uuid.UUID) ([]*Project, error)
}
