package main

import (
	"context"
	"log"

	"github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/config"
	"github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/server"
	"github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/storage"
)

func main() {
	// cfg - читает параметры конфигурации (хост, дб, миграция и дебаг)
	cfg := config.ReadConfig()
	log.Println(cfg)

	//stor - переменная для работы с хранилищем данных (БД)
	var stor server.Storage
	//инициализация хранилища (репозитория), передается контекст и строка подключения к БД
	stor, err := storage.NewRepo(context.Background(), cfg.DbDSN)
	// проверяет, возникла ли ошибка при создании хранилища
	if err != nil {
		log.Fatal(err.Error())
	}
	//
	if err = storage.Migrations(cfg.DbDSN, cfg.MigratePath); err != nil {
		log.Fatal(err.Error())
	}

	//
	server := server.New(cfg.Host, stor)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
