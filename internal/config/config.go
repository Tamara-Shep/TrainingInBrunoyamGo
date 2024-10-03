package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Host        string
	DbDSN       string
	MigratePath string
	Debug       bool
}

const (
	defaultDbDSN       = "postgres://postgres:6406655@localhost:5432/TRAINNINGINBRUNOYAMGO?sslmode=disable"
	defaultMigratePath = "maigrations"
	defaultHost        = ":8081"
)

// Ф-ия, которая конфигурирует наш проект
func ReadConfig() Config {
	var host string
	var dbDSN string
	var migratePath string

	flag.StringVar(&host, "host", defaultHost, "server host")
	flag.StringVar(&dbDSN, "db", defaultDbDSN, "data base addres")
	flag.StringVar(&migratePath, "m", defaultMigratePath, "path to migrations")

	debug := flag.Bool("debug", false, "enable debug loggin level")

	flag.Parse()

	hostEnv := os.Getenv("SERVER_HOST")
	dbDsnEnv := os.Getenv("DB_DSN")
	migratePathEnv := os.Getenv("MIGRATE_PATH")

	log.Println(hostEnv)

	if hostEnv != "" && host == defaultHost {
		host = hostEnv
	}
	if dbDsnEnv != "" && dbDSN == defaultDbDSN {
		dbDSN = dbDsnEnv
	}
	if migratePathEnv != "" && migratePath == defaultMigratePath {
		migratePath = migratePathEnv
	}

	return Config{
		Host:        host,
		DbDSN:       dbDSN,
		MigratePath: migratePath,
		Debug:       *debug,
	}
}
