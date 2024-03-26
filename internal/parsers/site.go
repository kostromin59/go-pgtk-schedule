package parsers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
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
		log.Println(err)
		return errors.New("ошибка получения сайта")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return errors.New("ошибка чтения сайта")
	}

	html := string(body)
	s.html = html

	return nil
}

func (s *Site) ExtractStudyYearId() (string, error) {
	if s.html == "" {
		return "", errors.New("html пустой")
	}

	re, err := regexp.Compile(`studyyear_id\s*:\s*'(\d+)'`)
	if err != nil {
		log.Println(err)
		return "", errors.New("ошибка компиляции регулярного выражения")
	}

	match := re.FindStringSubmatch(s.html)

	if len(match) <= 1 {
		return "", errors.New("совпадения не найдены")
	}

	return match[1], nil
}
