package schedule

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

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
		go func() {
			defer wg.Done()

			s.mu.Lock()
			defer s.mu.Unlock()

			b := scheduleBody{
				StudyYearId: studyYearId,
				Semester:    semester,
				StartDate:   week.StartDate.Initial,
				EndDate:     week.EndDate.Initial,
				GroupId:     group.Value,
			}

			lessons, err := parseByGroup(b)
			if err != nil {
				log.Println(err)
				return
			}

			s.Lessons[group.Value] = lessons
		}()
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
		log.Println(err)
		return nil, errors.New("[lesson, parseByGroup] ошибка маршализации")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
		return nil, errors.New("[lesson, parseByGroup] ошибка создания запроса")
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, errors.New("[lesson, parseByGroup] ошибка получения расписания")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("[lesson, parseByGroup] статус код не равен 200")
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, errors.New("[lesson, parseByGroup] ошибка чтения ответа")
	}

	var schedule []Lesson
	if err := json.Unmarshal(res, &schedule); err != nil {
		log.Println(err)
		return nil, errors.New("[lesson, parseByGroup] ошибка парсинга ответа")
	}

	return schedule, nil
}
