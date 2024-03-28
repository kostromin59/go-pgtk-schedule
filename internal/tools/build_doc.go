package tools

import (
	"errors"
	"io"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func BuildDoc(r io.Reader) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Println(err)
		return nil, errors.New("[site, buildDoc] ошибка создания документа")
	}

	return doc, nil
}
