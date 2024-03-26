package main

import (
	"fmt"

	"github.com/kostrominoff/go-pgtk-schedule/internal/schedule"
)

func main() {
	schedule := schedule.NewSchedule()
	schedule.Parse()
	fmt.Printf("%+v\n", schedule.Weekdates.Weeks)
}
