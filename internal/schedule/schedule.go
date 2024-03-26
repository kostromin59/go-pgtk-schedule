package schedule

import (
	"log"

	"github.com/kostrominoff/go-pgtk-schedule/internal/parsers"
)

type Schedule struct {
	Site      *parsers.Site
	Weekdates *parsers.Weekdates
}

func NewSchedule() *Schedule {
	return &Schedule{
		parsers.NewSite(),
		parsers.NewWeekdates(),
	}
}

func (s *Schedule) Parse() {
	if err := s.Site.Parse(); err != nil {
		log.Println(err)
		return
	}

	studyYearId, err := s.Site.ExtractStudyYearId()
	if err != nil {
		log.Println(err)
		return
	}

	if err := s.Weekdates.Parse(studyYearId); err != nil {
		log.Println(err)
		return
	}
}
