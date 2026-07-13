package infrastructure

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"app/core/project/domain"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(p *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: p}
}

func (r *PostgresRepository) Save(ctx context.Context, p *domain.Project) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO projects (id, name, owner_id, created_at) VALUES ($1, $2, $3, $4)`,
		p.ID, p.Name, p.OwnerID, p.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, owner uuid.UUID) ([]*domain.Project, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, owner_id, created_at FROM projects WHERE owner_id = $1 ORDER BY created_at DESC`,
		owner,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Project
	for rows.Next() {
		p := &domain.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.OwnerID, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
