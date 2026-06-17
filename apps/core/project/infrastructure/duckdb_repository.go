package infrastructure

import (
	"context"

	"github.com/google/uuid"

	"app/core/pkg/db"
	"app/core/project/domain"
)

type DuckDBRepository struct{ exec db.Executor }

func NewDuckDBRepository(e db.Executor) *DuckDBRepository { return &DuckDBRepository{exec: e} }

func (r *DuckDBRepository) Save(ctx context.Context, p *domain.Project) error {
	_, err := r.exec.ExecContext(ctx,
		`INSERT INTO projects (id, name, owner_id, created_at) VALUES (?, ?, ?, ?)`,
		p.ID.String(), p.Name, p.OwnerID.String(), p.CreatedAt,
	)
	return err
}

func (r *DuckDBRepository) List(ctx context.Context, owner uuid.UUID) ([]*domain.Project, error) {
	rows, err := r.exec.QueryContext(ctx,
		`SELECT id, name, owner_id, created_at FROM projects WHERE owner_id = ? ORDER BY created_at DESC`,
		owner.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Project
	for rows.Next() {
		var id, name, ownerID string
		p := &domain.Project{}
		if err := rows.Scan(&id, &name, &ownerID, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.ID, _ = uuid.Parse(id)
		p.OwnerID, _ = uuid.Parse(ownerID)
		p.Name = name
		out = append(out, p)
	}
	return out, rows.Err()
}
