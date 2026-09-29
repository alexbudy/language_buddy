package store

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
	_ "modernc.org/sqlite"
)

// Store wraps a SQLite database connection.
type Store struct {
	db *sql.DB
}

// Open opens (and if necessary creates) the SQLite database at path with
// foreign key enforcement turned on.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("store: open %s: %w", path, err)
	}

	if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("store: ping %s: %w", path, err)
    }

	
	return &Store{db: db}, nil
}

// Close closes the underlying SQLite database connection.
func (s *Store) Close() error { return s.db.Close() }


//go:embed schemas/main.sql
var mainSchema string

// EnsureSchema creates the database schema if it does not already exist.
func (s *Store) EnsureSchema(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, mainSchema); err != nil {
		return fmt.Errorf("store: ensure schema: %w", err)
	}
	log.Debug("Schema ensured: (%d bytes)", len(mainSchema))
	
	
	if err := s.migrate(ctx); err != nil {
	}
	
	return nil

}



func (s *Store) migrate(ctx context.Context) error {
	// add migrations here
	log.Debug("Running migrations")
	return nil
}