package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

type Config struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         int
}

func (c *Config) ToDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		c.Host,
		c.User,
		c.Password,
		c.DatabaseName,
		c.Port,
	)
}

func NewDatabase(c *Config) (*Database, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: c.ToDSN(),
	}))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	return &Database{db}, nil
}
