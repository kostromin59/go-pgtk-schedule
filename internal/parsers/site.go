package parsers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/kostrominoff/go-pgtk-schedule/internal/groups"
	"github.com/kostrominoff/go-pgtk-schedule/internal/tools"
)

type Site struct {
	html string
	mu   sync.Mutex
}

func NewSite() *Site {
	return &Site{}
}

func (s *Site) Parse() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	const url = "https://psi.thinkery.ru/shedule/public/public_shedule"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка получения сайта: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("статус код не равен 200: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения сайта: %w", err)
	}

	html := string(body)
	s.html = html

	return nil
}

func (s *Site) ExtractStudyYearId() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.html == "" {
		return "", errors.New("html пустой")
	}

	re, err := regexp.Compile(`studyyear_id\s*:\s*'(\d+)'`)
	if err != nil {
		return "", fmt.Errorf("ошибка компиляции регулярного выражения: %w", err)
	}

	match := re.FindStringSubmatch(s.html)

	if len(match) <= 1 {
		return "", errors.New("совпадения не найдены")
	}

	return match[1], nil
}

func (s *Site) ExtractSemester() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, err := tools.BuildDoc(bytes.NewBuffer([]byte(s.html)))
	if err != nil {
		return "", err
	}

	container := doc.Find("#termdiv")
	if container == nil {
		return "", errors.New("контейнер не найден")
	}

	option := container.Find("option").Last()
	if option == nil {
		return "", errors.New("выбранный семестр не найден")
	}

	semester, ok := option.Attr("value")
	if !ok {
		return "", errors.New("атрибут не найден")
	}

	return semester, nil
}

func (s *Site) ExtractGroups() ([]*groups.Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, err := tools.BuildDoc(bytes.NewBuffer([]byte(s.html)))
	if err != nil {
		return nil, err
	}

	container := doc.Find("#stream_iddiv")

	container.Find("option")

	var extractedGroups []*groups.Group

	const placeholder = "Выберите поток"
	container.Find("option").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		if name == placeholder {
			return
		}

		value, ok := s.Attr("value")
		if !ok {
			return
		}

		extractedGroups = append(extractedGroups, groups.NewGroup(name, value))
	})

	return extractedGroups, nil
}
