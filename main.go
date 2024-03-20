package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
) 

func main() {
	fmt.Println("Executing main.go...")

	scrapper := colly.NewCollector()
	var url_to_scrap = "https://es.wikipedia.org/wiki/Go_(lenguaje_de_programaci%C3%B3n)"

	scrapper.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting page: ", r.URL.String())
	})

	scrapper.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from page: ", r.Request.URL)
	})

	scrapper.OnError(func(_ * colly.Response, e error) {
		fmt.Println("An error occurred: ", e)
	})

	scrapper.OnHTML("main", func(e *colly.HTMLElement) {
		articleName := e.ChildText("span.mw-page-title-main")
		articleName = strings.TrimSpace(articleName)

		fmt.Println("Scrapped article title: ", articleName)
	})

	err := scrapper.Visit(url_to_scrap)
	
	if err != nil {
		log.Fatal(err)
	}
}

