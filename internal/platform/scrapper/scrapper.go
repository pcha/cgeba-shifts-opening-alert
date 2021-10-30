package scrapper

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Scrape(url, selector string) (*goquery.Selection, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	sel := doc.Find(selector)
	return sel, nil
}
