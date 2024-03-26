package parsers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
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

	const url = "https://psi.thinkery.ru/shedule/public/get_weekdates"

	type body struct {
		StudyYearId string `json:"studyyear_id"`
	}

	data := body{studyYearId}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return errors.New("[weekdates, Parse] ошибка маршализации")
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
		return errors.New("[weekdates, Parse] ошибка создания запроса")
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("[weekdates, Parse] ошибка получения дат")
	}

	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return errors.New("[weekdates, Parse] ошибка чтения ответа")
	}

	var weeks []Week
	if err := json.Unmarshal(res, &weeks); err != nil {
		log.Println(err)
		return errors.New("[weekdates, Parse] ошибка парсинга ответа")
	}

	w.Weeks = weeks

	return nil
}
