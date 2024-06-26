package schedule

import (
	"log/slog"
	"sync"

	"github.com/kostrominoff/go-pgtk-schedule/internal/groups"
	"github.com/kostrominoff/go-pgtk-schedule/internal/parsers"
)

type Schedule struct {
	Site      *parsers.Site
	Weekdates *parsers.Weekdates
	Groups    []*groups.Group
	Semester  string

	mu      sync.RWMutex
	Lessons map[string][]Lesson
}

type ScheduleByDates map[string][]Lesson

func NewSchedule() *Schedule {
	return &Schedule{
		Site:      parsers.NewSite(),
		Weekdates: parsers.NewWeekdates(),
		Lessons:   make(map[string][]Lesson),
	}
}

func (s *Schedule) Parse() {
	// Парсинг сайта
	if err := s.Site.Parse(); err != nil {
		slog.Error(err.Error())
		return
	}

	// Получение текущего года
	studyYearId, err := s.Site.ExtractStudyYearId()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Парсинг недель
	if err := s.Weekdates.Parse(studyYearId); err != nil {
		slog.Error(err.Error())
		return
	}

	// Получение семестра
	semester, err := s.Site.ExtractSemester()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	s.Semester = semester

	// Получение групп
	g, err := s.Site.ExtractGroups()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	s.Groups = g

	week := s.Weekdates.CurrentWeek()

	// Получение подгрупп
	var wg sync.WaitGroup
	for _, group := range s.Groups {
		wg.Add(1)
		go func(g *groups.Group) {
			defer wg.Done()
			if err := g.ParseSubgroups(studyYearId, semester, week.Value); err != nil {
				slog.Error(err.Error())
			}
		}(group)
	}

	wg.Wait()

	// Получение расписания
	s.ParseSchedules(studyYearId, semester, week)
}

func (s *Schedule) FindByGroup(groupId, subgroup string) ScheduleByDates {
	lessons := make(ScheduleByDates)

	for _, v := range s.Lessons[groupId] {
		key := v.DateStartText

		if _, ok := lessons[key]; !ok {
			lessons[key] = []Lesson{}
		}

		lessons[key] = append(lessons[key], v)
	}

	return lessons
}
