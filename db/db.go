package db

import (
	"campaignservice/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func LoadDBConfig() models.DBConfig {
	return models.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

func Connect() *sql.DB {
	cfg := LoadDBConfig()
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}
	return db
}
