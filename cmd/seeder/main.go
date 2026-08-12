package main

import (
	"log"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/seeder"
)

func main() {
	cfg := config.LoadConfig()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatal("failed to connect db: ", err)
	}

	if err := seeder.SeedAll(db); err != nil {
		log.Fatal("seeding failed: ", err)
	}
}
