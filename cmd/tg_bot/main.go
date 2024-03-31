package main

import (
	"fmt"

	"github.com/kostrominoff/go-pgtk-schedule/internal/schedule"
)

func main() {
	schedule := schedule.NewSchedule()
	schedule.Parse()
	s := schedule.FindByGroup("1925", "")
	fmt.Printf("%#v\n", s)
	// schedule.Message(schedule.Lessons["1925"])
}
