package schedule

import (
	"fmt"
	"log"

	"github.com/kostrominoff/go-pgtk-schedule/internal/groups"
	"github.com/kostrominoff/go-pgtk-schedule/internal/parsers"
)

type Schedule struct {
	Site      *parsers.Site
	Weekdates *parsers.Weekdates
	Groups    []*groups.Group
	Semester  string
}

func NewSchedule() *Schedule {
	return &Schedule{
		Site:      parsers.NewSite(),
		Weekdates: parsers.NewWeekdates(),
	}
}

func (s *Schedule) Parse() {
	// Парсинг сайта
	if err := s.Site.Parse(); err != nil {
		log.Println(err)
		return
	}

	// Получение текущего года
	studyYearId, err := s.Site.ExtractStudyYearId()
	if err != nil {
		log.Println(err)
		return
	}

	// Парсинг недель
	if err := s.Weekdates.Parse(studyYearId); err != nil {
		log.Println(err)
		return
	}

	// Получение семестра
	semester, err := s.Site.ExtractSemester()
	if err != nil {
		log.Println(err)
		return
	}

	s.Semester = semester
	fmt.Println(semester)

	// Получение групп
	groups, err := s.Site.ExtractGroups()
	if err != nil {
		log.Println(err)
		return
	}

	s.Groups = groups

	// Получение подгрупп
	for _, group := range s.Groups {
		// log.Println(group)
		if err := group.ParseSubgroups(); err != nil {
			log.Println(err)
		}
	}
}
