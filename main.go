package main

import (
	"context"
	"flag"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
)

const defaultDBPath = "my_db.db" // sqLite db path

var debug bool
var debugSQLStmts bool

func main() {
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.BoolVar(&debugSQLStmts, "debug-sql", false, "enable SQL debug logging")
	flag.Parse()

	log.DebugEnabled = debug
	log.DebugSQLEnabled = debugSQLStmts

	log.Debug("Starting up...")
	err := runInit()
	if err != nil {panic(err)}
}

func runInit()  error {
	s, err := store.Open(defaultDBPath)
	if err != nil {
		return err
	}

	defer s.Close()

	log.Debug("Database ready at %s\n", defaultDBPath)

	ctx := context.Background()

	// Ensure schema is loaded
	if err := s.EnsureSchema(ctx); err != nil {
		return err
	}

	// Seed words from csv
	if err := s.SeedWords(ctx); err != nil {
		return err
	}



	return nil
}