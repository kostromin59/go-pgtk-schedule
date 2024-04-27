package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/kostrominoff/go-pgtk-schedule/internal/tools"
)

type Week struct {
	Value     int              `json:"value"`
	StartDate tools.CustomDate `json:"start_date"`
	EndDate   tools.CustomDate `json:"end_date"`
	Selected  bool             `json:"selected"`
}

type Weekdates struct {
	Weeks []Week
	mu    sync.Mutex
}

func NewWeekdates() *Weekdates {
	return &Weekdates{}
}

func (w *Weekdates) Parse(studyYearId string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	const url = "https://psi.thinkery.ru/shedule/public/get_weekdates_actual"

	type body struct {
		StudyYearId string `json:"studyyear_id"`
	}

	data := body{studyYearId}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка маршализации: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка получения дат: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("статус код не равен 200: %v", resp.StatusCode)
	}

	var weeks []Week

	if err := json.NewDecoder(resp.Body).Decode(&weeks); err != nil {
		return fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	w.Weeks = weeks

	return nil
}

func (w *Weekdates) CurrentWeek() *Week {
	var selected int

	for i, v := range w.Weeks {
		if v.Selected {
			selected = i
			break
		}
	}

	// TODO: Расскоментировать
	// now := time.Now()

	// if (now.Weekday() == time.Saturday && now.Hour() >= 12) || now.Weekday() == time.Sunday {
	// selected++
	// }

	return &w.Weeks[selected]
}
