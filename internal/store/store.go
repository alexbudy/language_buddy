package store

import (
	"database/sql"
	"fmt"

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

