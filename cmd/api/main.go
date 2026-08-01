package main

import (
	"log"

	"github.com/bachtiarrizaa/sembako-be/internal/bootstrap"
	"github.com/bachtiarrizaa/sembako-be/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	app, err := bootstrap.InitializeApp(cfg)
	if err != nil {
		log.Fatal("failed to initialize app: ", err)
	}

	log.Printf("Server running on port %s", cfg.AppPort)
	if err := app.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("failed to run server: ", err)
	}
}
