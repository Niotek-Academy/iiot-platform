package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store bundles the live connection pool with the generated query methods.
type Store struct {
	Pool *pgxpool.Pool
	*generated.Queries
}

// Connect opens a pooled connection to Postgres and verifies it with a ping.
// It retries a few times because the Postgres container may still be
// starting up when the backend boots (common in docker-compose setups).
func Connect(ctx context.Context, databaseURL string) (*Store, error) {
	var pool *pgxpool.Pool
	var err error

	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		pool, err = pgxpool.New(ctx, databaseURL)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			err = pool.Ping(pingCtx)
			cancel()
			if err == nil {
				break
			}
		}

		if attempt == maxAttempts {
			return nil, fmt.Errorf("could not connect to database after %d attempts: %w", maxAttempts, err)
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return &Store{
		Pool:    pool,
		Queries: generated.New(pool),
	}, nil
}

// Close releases the connection pool. Call this on graceful shutdown.
func (s *Store) Close() {
	if s.Pool != nil {
		s.Pool.Close()
	}
}
