package main

import (
	"flag"
	"log"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/seeder"
)

func main() {
	name := flag.String("name", "", "Specify seeder name (all, roles, permissions, users, store-config, loyalty, demo)")
	target := flag.String("target", "", "Alias for -name")
	flag.Parse()

	seederName := *name
	if seederName == "" {
		seederName = *target
	}
	if seederName == "" {
		seederName = "all"
	}

	cfg := config.LoadConfig()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatal("failed to connect db: ", err)
	}

	if err := seeder.SeedByName(db, seederName); err != nil {
		log.Fatal("seeding failed: ", err)
	}

	log.Printf("Seeder '%s' executed successfully\n", seederName)
}
