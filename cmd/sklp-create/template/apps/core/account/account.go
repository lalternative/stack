// Package account is the core's account/auth context. It is intentionally
// minimal: authentication itself is owned by the web app (better-auth), and the
// core only verifies the minted JWT (see apps/core/middleware). What lives here
// is the small set of account-facing actions the core still needs to serve —
// today, clearing the session cookie on logout.
//
// It is a plain service (no CQRS/event-sourcing like the example context) on
// purpose: these are stateless auth-cookie operations, not an event-sourced
// aggregate. Grow it into the DDD layout only if account state moves into the
// core.
package account

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service is the account context's HTTP-facing facade. It holds the pool so
// future account queries (e.g. reading a profile the core owns) have a home.
type Service struct {
	pool *pgxpool.Pool
}

// NewService wires the account service to the shared pgx pool.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}
