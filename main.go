package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	var city string
	fmt.Println("Enter a city name:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	city = scanner.Text()
	url := "https://en.wikipedia.org/wiki/" + city
	rows := scraper(url)
	var state, population, area string
	var re = regexp.MustCompile(`Total|Metro|• Metro|• Total`)
	for i := 0; i < len(rows); i++ {
		switch {
		case strings.Contains(rows[i].label, "State"):
			state = string(rows[i].data)
		case re.MatchString(rows[i].label):
			if strings.Contains(rows[i].data, "km") {
				area = rows[i].data
			} else if isNumber(rows[i].data) {
				population = rows[i].data
			}
		}
	}
	fmt.Println("State: " + state)
	fmt.Println("Population: " + population)
	fmt.Println("Area: " + area)
	duration := time.Duration(60) * time.Second
	time.Sleep(duration)
}

func scraper(url string) []mergedRow {
	var rows []mergedRow
	collector := colly.NewCollector()
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got response from", r.Request.URL)
	})
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	collector.OnHTML(".mergedrow", func(r *colly.HTMLElement) {
		mer := mergedRow{}
		var re = regexp.MustCompile(`\([^()]*\)`)
		mer.label = re.ReplaceAllString(r.ChildText(".infobox-label"), "$1")
		mer.data = re.ReplaceAllString(r.ChildText(".infobox-data"), "$1")
		rows = append(rows, mer)
	})
	collector.Visit(url)
	return rows
}

func isNumber(str string) bool {
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		if runes[i] != 48 && runes[i] != 49 && runes[i] != 50 && runes[i] != 51 && runes[i] != 52 && runes[i] != 53 && runes[i] != 54 && runes[i] != 55 && runes[i] != 56 && runes[i] != 57 && runes[i] != 44 {
			return false
		}
	}
	return true
}

type mergedRow struct {
	label, data string
}
