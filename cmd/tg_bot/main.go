package main

import (
	"log"

	"github.com/kostrominoff/go-pgtk-schedule/internal/config"
	"github.com/kostrominoff/go-pgtk-schedule/internal/logger"
	"github.com/kostrominoff/go-pgtk-schedule/internal/schedule"
	"github.com/kostrominoff/go-pgtk-schedule/internal/storage"
)

func main() {
	logger.New()

	config := config.NewTgBot()
	storage, err := storage.NewDatabase(&storage.Config{
		Host:         config.DBHost,
		User:         config.DBUser,
		Password:     config.DBPassword,
		DatabaseName: config.DBName,
		Port:         config.DBPort,
	})
	if err != nil {
		log.Fatal(err)
	}

	storage.Migrate()

	schedule := schedule.NewSchedule()
	schedule.Parse()

	// s := schedule.FindByGroup("1925", "")
	// fmt.Printf("%#v\n", s)
	// schedule.Message(schedule.Lessons["1925"])
}
