package parsers

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/kostrominoff/go-pgtk-schedule/internal/groups"
	"golang.org/x/net/html"
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

func (s *Site) ExtractGroups() ([]*groups.Group, error) {
	doc, err := html.Parse(bytes.NewBuffer([]byte(s.html)))
	if err != nil {
		log.Println(err)
		return nil, errors.New("ошибка парсинга документа")
	}

	var result []*groups.Group

	var processGroupsContainer func(*html.Node)
	var extractOptions func(*html.Node)

	extractOptions = func(n *html.Node) {
		// Поиск option
		if n.Type == html.ElementNode && n.Data == "option" {
			name := extractText(n)

			// Поиск значения группы
			var value string
			for _, attr := range n.Attr {
				if attr.Key == "value" {
					value = attr.Val
				}
			}

			result = append(result, groups.NewGroup(name, value))
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractOptions(c)
		}
	}

	processGroupsContainer = func(n *html.Node) {
		// Поиск div
		if n.Type == html.ElementNode && n.Data == "div" {
			// Поиск аттрибута id со значением stream_iddiv
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "stream_iddiv" {
					extractOptions(n)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processGroupsContainer(c)
		}
	}

	processGroupsContainer(doc)

	return result, nil
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return text
}
