package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/kostrominoff/go-pgtk-schedule/internal/groups"
	"github.com/kostrominoff/go-pgtk-schedule/internal/parsers"
)

type Lesson struct {
	CabinetFullNumber string `json:"cabinet_fullnumber_wotype"`
	DisciplineName    string `json:"discipline_name"`
	DateStartText     string `json:"date_start_text"`
	DaytimeName       string `json:"daytime_name"`
	WeekdayName       string `json:"weekday_name"`
	TeacherFio        string `json:"teacher_fio"`
	StreamName        string `json:"stream_name"`
	SubgroupName      string `json:"subgroup_name"`
	ClassTypeName     string `json:"classtype_name"`
	GroupId           int    `json:"stream_id"`
	DaytimeStart      string `json:"daytime_start"`
}

func (s *Schedule) ParseSchedules(studyYearId string, semester string, week *parsers.Week) {
	var wg sync.WaitGroup
	for _, group := range s.Groups {
		wg.Add(1)
		go func(g *groups.Group) {
			defer wg.Done()

			s.mu.Lock()
			defer s.mu.Unlock()

			b := scheduleBody{
				StudyYearId: studyYearId,
				Semester:    semester,
				StartDate:   week.StartDate.Initial,
				EndDate:     week.EndDate.Initial,
				GroupId:     g.Value,
			}

			lessons, err := parseByGroup(b)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			s.Lessons[g.Value] = lessons
		}(group)
	}

	wg.Wait()
}

type scheduleBody struct {
	StudyYearId string `json:"studyyear_id"`
	GroupId     string `json:"stream_id"`
	Semester    string `json:"term"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

func parseByGroup(b scheduleBody) ([]Lesson, error) {
	const url = "https://psi.thinkery.ru/shedule/public/public_getsheduleclasses_spo"

	jsonBody, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршализации: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения расписания: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("статус код не равен 200: %v", resp.StatusCode)
	}

	var schedule []Lesson

	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	return schedule, nil
}
