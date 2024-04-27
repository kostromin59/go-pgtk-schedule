package tools

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

func BuildDoc(r io.Reader) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания документа: %w", err)
	}

	return doc, nil
}
