package main

import (
	"flag"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
)

const defaultDBPath = "my_db.db" // sqLite db path

var debug bool

func main() {
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()

	log.DebugEnabled = debug

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

	return nil
}