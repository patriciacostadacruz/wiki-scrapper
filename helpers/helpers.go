package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
)

type ArticleData struct {
  Name  string `json:"articleName"`
  Intro string `json:"articleIntro"`
}

type FakeUserAgentResponse struct {
    Result []string `json:"result"`
}

func GetRandomUserAgent(userAgentList []string) string {
    if len(userAgentList) == 0 {
        log.Fatal("User agent list is empty")
    }
    randomIndex := rand.Intn(len(userAgentList))
    return userAgentList[randomIndex]
}

func GetUserAgentList() []string {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("No .env file found")
    }

    scrapeopsAPIKey, ok := os.LookupEnv("SCRAPE_OPS_API_KEY")
    if !ok {
        log.Fatalf("Error loading data from .env file")
    }
    scrapeopsAPIEndpoint := "http://headers.scrapeops.io/v1/user-agents?api_key=" + scrapeopsAPIKey

    req, _ := http.NewRequest("GET", scrapeopsAPIEndpoint, nil)
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    res, err := client.Do(req)
    if err == nil && res.StatusCode == 200 {
        defer res.Body.Close()

        var fakeUserAgentResponse FakeUserAgentResponse
        json.NewDecoder(res.Body).Decode(&fakeUserAgentResponse)
        return fakeUserAgentResponse.Result
    }

    var emptySlice []string
    return emptySlice
}

func ScrapeArticle(e *colly.HTMLElement) ArticleData {
    article := ArticleData{}
    article.Name = e.ChildText("span.mw-page-title-main")
    article.Name = strings.TrimSpace(article.Name)
    // TODO: improve selector to pick intro only mf-section-0
    article.Intro = e.ChildText("div.mw-parser-output")
    article.Intro = strings.TrimSpace(article.Intro)

    // Removes references like [1], [2], etc
    re := regexp.MustCompile(`\[\d+]`)
    article.Intro = re.ReplaceAllString(article.Intro, "")

    return article
}

func WriteArticlesToFile(articles []ArticleData) {
    js, err := json.MarshalIndent(articles, "", "    ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("[INFO] Writing data to file...")
    if err := os.WriteFile("articles-data.json", js, 0664); err == nil {
        fmt.Println("[END] Data written to file successfully!")
    }
}
