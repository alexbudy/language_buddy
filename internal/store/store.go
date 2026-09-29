package store

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
	_ "modernc.org/sqlite"
)

// File for setting up the store

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

//go:embed data/words.csv
var wordsCSV string

// SeedWords seeds the database with the initial word data.
func (s *Store) SeedWords(ctx context.Context) error {
	log.Debug("Seeding words...")

	var count int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM words").Scan(&count); err != nil {
		return fmt.Errorf("store: count words: %w", err)
	}
	if count > 0 {
		log.Debug("%d words already seeded, not seeding", count)
		return nil
	}

	reader := csv.NewReader(strings.NewReader(wordsCSV))
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("store: read words header: %w", err)
	}
	log.Debug("Found header: %v", header)

	colIdx := make(map[string]int, len(header))
	for i, name := range header {
		colIdx[name] = i
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: seed words: %w", err)
	}
	defer tx.Rollback()

	instStmt := "INSERT INTO words (english, spanish, ukrainian) VALUES (?, ?, ?)"
	stmt, err := tx.PrepareContext(ctx, instStmt)
	if err != nil {
		return fmt.Errorf("store: seed words: %w", err)
	}
	defer stmt.Close()

	log.DebugSQL("Running statement for each word: %s", instStmt)

	for {
		row, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("store: read words row: %w", err)
		}
		english := row[colIdx["english"]]
		spanish := row[colIdx["spanish"]]
		ukrainian := row[colIdx["ukrainian"]]
		if _, err := stmt.ExecContext(ctx, english, spanish, ukrainian); err != nil {
			return fmt.Errorf("store: insert word %q: %w", english, err)
		}
		log.DebugSQL("Inserted word %q into words table", english)
	}

	return tx.Commit()
}

func (s *Store) migrate(ctx context.Context) error {
	// add migrations here
	log.Debug("Running migrations...")
	newMigrations := 0
	
	log.Debug("New migrations ran: %d", newMigrations)
	return nil
}

