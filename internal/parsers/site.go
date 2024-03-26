package parsers

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/kostrominoff/go-pgtk-schedule/internal/groups"
)

type Site struct {
	html string
	mu   sync.RWMutex
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
		log.Println(err)
		return errors.New("[site, Parse] ошибка получения сайта")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return errors.New("[site, Parse] ошибка чтения сайта")
	}

	html := string(body)
	s.html = html

	return nil
}

func (s *Site) ExtractStudyYearId() (string, error) {
	if s.html == "" {
		return "", errors.New("[site, ExtractStudyYearId] html пустой")
	}

	re, err := regexp.Compile(`studyyear_id\s*:\s*'(\d+)'`)
	if err != nil {
		log.Println(err)
		return "", errors.New("[site, ExtractStudyYearId] ошибка компиляции регулярного выражения")
	}

	match := re.FindStringSubmatch(s.html)

	if len(match) <= 1 {
		return "", errors.New("[site, ExtractStudyYearId] совпадения не найдены")
	}

	return match[1], nil
}

func (s *Site) buildDoc() (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(s.html)))
	if err != nil {
		log.Println(err)
		return nil, errors.New("[site, buildDoc] ошибка создания документа")
	}

	return doc, nil
}

func (s *Site) ExtractSemester() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, err := s.buildDoc()
	if err != nil {
		return "", err
	}

	container := doc.Find("#termdiv")
	if container == nil {
		return "", errors.New("[site, ExtractSemester] контейнер не найден")
	}

	option := container.Find("option").Last()
	if option == nil {
		return "", errors.New("[site, ExtractSemester] выбранный семестр не найден")
	}

	semester, ok := option.Attr("value")
	if !ok {
		return "", errors.New("[site, ExtractSemester] атрибут не найден")
	}

	return semester, nil
}

func (s *Site) ExtractGroups() ([]*groups.Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, err := s.buildDoc()
	if err != nil {
		return nil, err
	}

	container := doc.Find("#stream_iddiv")
	if container == nil {
		return nil, errors.New("[site, ExtractGroups] контейнер не найден")
	}

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
