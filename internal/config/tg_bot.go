package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type TgBot struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     int
}

func NewTgBot() *TgBot {
	config := &TgBot{}

	if err := godotenv.Load("configs/.env"); err != nil {
		log.Fatal(err)
	}

	host := os.Getenv("POSTGRES_HOST")

	if host == "" {
		slog.Warn("POSTGRES_HOST пустой")
		host = "localhost"
	}

	user := os.Getenv("POSTGRES_USER")

	if user == "" {
		log.Fatal("POSTGRES_USER пустой")
	}

	password := os.Getenv("POSTGRES_PASSWORD")

	if password == "" {
		log.Fatal("POSTGRES_PASSWORD пустой")
	}

	name := os.Getenv("POSTGRES_DB")

	if name == "" {
		log.Fatal("POSTGRES_DB пустой")
	}

	var port int = 5432
	portEnv := os.Getenv("POSTGRES_PORT")

	if portEnv == "" {
		slog.Warn("POSTGRES_PORT пустой")
	} else {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			log.Fatal()
		}

		port = portInt
	}

	config.DBHost = host
	config.DBUser = user
	config.DBPassword = password
	config.DBName = name
	config.DBPort = port

	return config
}
