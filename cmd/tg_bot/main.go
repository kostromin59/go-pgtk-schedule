package main

import (
	"log"

	"github.com/kostrominoff/go-pgtk-schedule/internal/parsers"
)

func main() {
	s := parsers.NewSite()
	if err := s.Parse(); err != nil {
		log.Fatal(err)
	}

	year, err := s.ExtractStudyYearId()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(year)
}
