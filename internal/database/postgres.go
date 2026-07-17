package database

import (
	"database/sql" // standard-library DB abstraction; driver-agnostic
	"fmt"
	"log/slog"
	"time"

	// Blank import: we never reference lib/pq directly, but importing it for its
	// side effects runs its init(), which registers the "postgres" driver with
	// database/sql. The leading _ tells the compiler "import but don't use" so it
	// doesn't error on the unused package.
	_ "github.com/lib/pq"
)

// NewPostgresConnection opens a connection pool to Postgres, tunes the pool, and
// verifies it's actually reachable before handing it back. Returning the error
// lets the caller decide how to fail (database/sql defers real connecting, so a
// bad config wouldn't surface without the explicit Ping below).
func NewPostgresConnection(connStr string) (*sql.DB, error) {
	// sql.Open does NOT establish a connection — it only validates arguments and
	// prepares a lazy *sql.DB pool. Actual connections are made on first use.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// %w wraps the original error so callers can still inspect it with
		// errors.Is / errors.As, while adding human-readable context here.
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// *sql.DB is a pool, not a single connection. These settings bound how many
	// connections it holds so we don't exhaust the database's connection limit
	// while still reusing connections for performance.
	db.SetMaxOpenConns(25)                 // hard cap on simultaneous open connections
	db.SetMaxIdleConns(10)                 // keep up to 10 warm/idle for quick reuse
	db.SetConnMaxLifetime(5 * time.Minute) // recycle conns periodically (avoids stale ones behind load balancers / server timeouts)
	db.SetConnMaxIdleTime(1 * time.Minute) // close idle conns after a minute to free resources

	// Ping forces one real round-trip to the server, so we fail fast at startup
	// if the DB is unreachable or credentials are wrong — rather than on the
	// first query much later.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("database connection established")
	return db, nil
}
