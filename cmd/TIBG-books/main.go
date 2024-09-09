package main

import (
	"log"

	"github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/config"
	"github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/server"
	"github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/storage"
)

func main() {
	cfg := config.ReadConfig()
	log.Print(cfg)
	storage := storage.New()
	server := server.New(cfg.Host, storage)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
