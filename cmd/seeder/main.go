package main

import (
	"log"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

func main() {
	cfg := config.LoadConfig()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatal("failed to connect db: ", err)
	}

	roles := []entity.Role{
		{Name: "Owner"},
		{Name: "Kasir"},
	}

	for _, role := range roles {
		result := db.FirstOrCreate(&role, entity.Role{Name: role.Name})
		if result.Error != nil {
			log.Fatal("failed to seed role: ", result.Error)
		}
	}

	log.Println("seeding roles done!")
}
