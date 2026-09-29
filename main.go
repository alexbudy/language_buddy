package main

import (
	"fmt"

	"github.com/alexbudy/go_spanish_rewrite/internal/store"
)

const defaultDBPath = "my_db.db" // sqLite db path

func main() {
	fmt.Println("Hello World")
	err := runInit()
	if err != nil {panic(err)}
}

func runInit()  error {
	s, err := store.Open(defaultDBPath)
	
	defer s.Close()

	fmt.Println("Hello world opened")

	return err
}