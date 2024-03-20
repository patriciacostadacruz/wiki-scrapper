package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
) 

func main() {
	fmt.Println("Executing main.go...")
	var articleName string

	scrapper := colly.NewCollector()
	
	var url_to_scrap = "https://es.wikipedia.org/wiki/Go_(lenguaje_de_programaci%C3%B3n)"

	scrapper.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting page: ", r.URL.String())
		// r.Headers.Set("User-Agent", "to_be_created")
	})

	scrapper.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from page: ", r.Request.URL)
	})

	
	scrapper.OnHTML("main", func(e *colly.HTMLElement) {
		articleName = e.ChildText("span.mw-page-title-main")
		articleName = strings.TrimSpace(articleName)
	})
	
	scrapper.OnError(func(r * colly.Response, e error) {
		log.Fatal("An error occurred: ", e)
	})

	scrapper.OnScraped(func(r *colly.Response) {
		fmt.Println("Scrapped article title: ", articleName)
		js, err := json.MarshalIndent(articleName, "", "    ")
      if err != nil {
        log.Fatal(err)
      }
      fmt.Println("Writing data to file...")
      if err := os.WriteFile("article-data.json", js, 0664); err == nil {
          fmt.Println("Data written to file successfully.")
      }
	})

	scrapper.Visit(url_to_scrap)
}

