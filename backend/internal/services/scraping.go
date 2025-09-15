package services

import (
	"strings"

	"github.com/gocolly/colly"
)

type ScrapingService interface {
	Scrape(url string) (string, error)
}

type scrapingService struct{}

func NewScraper() *scrapingService {
	return &scrapingService{}
}

func (s *scrapingService) Scrape(url string) (string, error) {
	c := colly.NewCollector()
	var content strings.Builder

	c.OnHTML("body", func(e *colly.HTMLElement) {
		content.WriteString(e.Text)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return content.String(), nil
}
