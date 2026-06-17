// Package db exposes a thin Executor interface that lets repositories
// stay backend-agnostic. The default backend is DuckDB; swap by
// implementing Executor with another driver (pgx, sqlite, ...).
package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/marcboeker/go-duckdb"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func Open(path string) (Executor, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return sql.Open("duckdb", path)
}

// Migrate applies every *.sql file in dir in lexicographic order. The
// pattern matches `apps/core/migrations/` in the skalpai core service.
func Migrate(ctx context.Context, exec Executor, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, n := range names {
		b, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			return fmt.Errorf("read %s: %w", n, err)
		}
		if _, err := exec.ExecContext(ctx, string(b)); err != nil {
			return fmt.Errorf("apply %s: %w", n, err)
		}
	}
	return nil
}
