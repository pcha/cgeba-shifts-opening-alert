package shifttable

import (
	"cgeba-shift-opening-alerter/internal/platform/scrapper"

	"github.com/PuerkitoBio/goquery"
)

const tableURL = "https://www.cgeonline.com.ar/informacion/apertura-de-citas.html"
const tableSelector = ".table-bordered tbody tr"

type TableRecord struct {
	LastOpening string
	NextOpening string
}

type ShiftTable map[string]TableRecord

func GetTable() (ShiftTable, error) {
	sel, err := scrapper.Scrape(tableURL, tableSelector)
	if err != nil {
		return nil, err
	}
	table := parseTable(sel)
	return table, nil
}

func parseTable(selection *goquery.Selection) ShiftTable {
	table := ShiftTable{}
	selection.Each(func(i int, record *goquery.Selection) {
		cells := record.Find("td")
		section := cells.Eq(0).Text()
		lastOpening := cells.Eq(1).Text()
		nextOpening := cells.Eq(2).Text()
		table[section] = TableRecord{
			LastOpening: lastOpening,
			NextOpening: nextOpening,
		}
	})
	return table
}
