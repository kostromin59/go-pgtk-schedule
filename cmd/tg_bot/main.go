package main

import (
	"github.com/kostrominoff/go-pgtk-schedule/internal/schedule"
)

func main() {
	schedule := schedule.NewSchedule()
	schedule.Parse()
}
