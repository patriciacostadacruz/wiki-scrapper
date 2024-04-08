package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/patriciacostadacruz/wiki-scrapper/helpers"

	"github.com/gocolly/colly/v2"
)

var urlsToScrap = []string{
    "https://en.wikipedia.org/wiki/Go_(programming_language)",
    "https://en.wikipedia.org/wiki/Python",
    "https://en.wikipedia.org/wiki/Solidity",
}

func main() {
    fmt.Println("[STARTING] Executing program...")

    scrapper := colly.NewCollector()

    var articles []helpers.ArticleData

    scrapper.Limit(&colly.LimitRule{
        Delay:       1,
        RandomDelay: time.Duration(rand.Intn(10)),
    })

    userAgentList := helpers.GetUserAgentList()

    scrapper.OnRequest(func(r *colly.Request) {
        fmt.Println("[INFO] Requesting page: ", r.URL.String())
        r.Headers.Set("User-Agent", helpers.GetRandomUserAgent(userAgentList))
    })

    scrapper.OnResponse(func(r *colly.Response) {
        fmt.Println("[INFO] Got a response from page: ", r.Request.URL)
    })

    scrapper.OnHTML("main", func(e *colly.HTMLElement) {
        article := helpers.ScrapeArticle(e)
        articles = append(articles, article)
    })

    scrapper.OnError(func(r *colly.Response, e error) {
        log.Fatal("An error occurred: ", e)
    })

    for _, page := range urlsToScrap {
        scrapper.Visit(page)
    }

    helpers.WriteArticlesToFile(articles)
}
