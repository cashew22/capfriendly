package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	baseURL  = "https://www.capfriendly.com/browse/active/"
	year     = "2024"
	order    = "/aav"
	filter   = "?stats-season=" + year + "&display=expiry-year,aav&hide=clauses,handed,expiry-status,salary,caphit,skater-stats,goalie-stats"
	page     = "&pg="
	url      = baseURL + year + order + filter + page
	lastPage = 31
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("capfriendly.com", "www.capfriendly.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	var headers []string
	c.OnHTML("th", func(e *colly.HTMLElement) {
		headers = append(headers, strings.TrimSpace(e.Text))
		c.OnHTMLDetach("th")
	})

	var players [][]string
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		var player []string
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			player = append(player, strings.TrimSpace(el.Text))
		})
		if len(player) > 0 {
			player[0] = strings.SplitN(player[0], " ", 2)[1]
			players = append(players, player)
		}
	})

	page := 1
	for page <= lastPage {
		c.Visit(url + fmt.Sprintf("%d", page))
		page++
	}

	saveCSV(headers, players, "capfriendly.csv")
}

func saveCSV(headers []string, players [][]string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(headers)
	for _, player := range players {
		writer.Write(player)
	}
}
