package scraping

import (
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/jgabriel1/mh-weakness-bot/lib/element"
	"github.com/jgabriel1/mh-weakness-bot/lib/hitzones"
)

const baseScrapeURL = "https://mhworld.kiranico.com"

func ScrapeMonsterHitzonesTable(monsterPath string) (*hitzones.Table, error) {
	t := hitzones.NewTable()
	c := colly.NewCollector()

	c.OnHTML("table", func(tableElement *colly.HTMLElement) {
		if !isPhysiologyTable(tableElement) {
			return
		}

		lookupIndexes := make(map[int]int)
		tableElement.DOM.Find("thead > tr > th").Each(func(_ int, s *goquery.Selection) {
			iconSrc, srcExists := findIconSrcInTableHeader(s)
			if !srcExists {
				return
			}

			el := findElementByIconSrc(iconSrc)
			if el == element.Unknown {
				return
			}

			colIndex := t.AddColumn(el)
			lookupIndexes[colIndex] = s.Index()
		})

		tableElement.DOM.Find("tbody > tr").Each(func(_ int, sr *goquery.Selection) {
			for colIndex, lookupIndex := range lookupIndexes {
				value, err := findTableRowValueForColumnIndex(sr, lookupIndex)
				if err == nil {
					t.AddValueToColumn(colIndex, value)
				}
			}
		})
	})

	if err := c.Visit(baseScrapeURL + monsterPath); err != nil {
		return nil, err
	}

	return t, nil
}

func isPhysiologyTable(e *colly.HTMLElement) bool {
	title := e.DOM.Parent().Prev()
	return title.Is("h6") && title.Text() == "Physiology"
}

func findIconSrcInTableHeader(s *goquery.Selection) (string, bool) {
	iconSrc, srcExists := s.Children().Attr("src")
	return iconSrc, s.Children().Is("img") && srcExists
}

func findElementByIconSrc(imgSrc string) element.Element {
	re := regexp.MustCompile(`element_(\d+)\.png`)
	matches := re.FindStringSubmatch(imgSrc)

	var elementNumber string
	if len(matches) > 1 {
		elementNumber = matches[1]
	}

	switch elementNumber {
	case "1":
		return element.Fire
	case "2":
		return element.Water
	case "3":
		return element.Ice
	case "4":
		return element.Thunder
	case "5":
		return element.Dragon
	default:
		return element.Unknown
	}
}

func findTableRowValueForColumnIndex(s *goquery.Selection, index int) (int, error) {
	return strconv.Atoi(s.Find("td").Get(index).FirstChild.Data)
}
