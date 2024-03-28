package groups

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/kostrominoff/go-pgtk-schedule/internal/tools"
)

type Group struct {
	Name      string
	Value     string
	Subgroups []string
	mu        sync.Mutex
}

func NewGroup(name, value string) *Group {
	return &Group{Name: name, Value: value}
}

func (g *Group) ParseSubgroups(studyYearId, semester string, weekNumber int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	const url = "https://psi.thinkery.ru/shedule/public/public_shedule_spo_grid"

	type body struct {
		StudyYearId string `json:"studyyear_id"`
		StreamId    string `json:"stream_id"`
		Term        string `json:"term"`
		Dateweek    int    `json:"dateweek"`
	}

	data := body{
		StudyYearId: studyYearId,
		StreamId:    g.Value,
		Term:        semester,
		Dateweek:    weekNumber,
	}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return errors.New("[groups, ParseSubgroups] ошибка маршализации")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
		return errors.New("[groups, ParseSubgroups] ошибка создания запроса")
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("[groups, ParseSubgroups] ошибка получения подгрупп")
	}

	defer resp.Body.Close()

	subgroups, err := extractSubgroups(resp.Body)
	if err != nil {
		log.Println(err)
		return errors.New("[groups, ParseSubgroups] ошибка извлечения подгрупп")
	}

	g.Subgroups = subgroups

	return nil
}

func extractSubgroups(r io.Reader) ([]string, error) {
	doc, err := tools.BuildDoc(r)
	if err != nil {
		return nil, err
	}

	table := doc.Find("#timetable")
	head := table.Find("thead")
	tableRow := head.Find("tr")

	if tableRow.Length() < 2 {
		return nil, errors.New("подгруппы не найдены")
	}

	lastRow := tableRow.Last()

	subgroups := make([]string, 0, 3)
	lastRow.Find("th").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text == "" {
			return
		}

		subgroups = append(subgroups, text)
	})

	return subgroups, nil
}
