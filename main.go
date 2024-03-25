package main

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

	"github.com/joho/godotenv"

	"github.com/gocolly/colly/v2"
)

type ArticleData struct {
	Name  string `json:"articleName"`
	Intro string `json:"articleIntro"`
}

type FakeUserAgentResponse struct {
	Result				[]string `json:"result"`
}

var urls_to_scrap = []string{
		"https://en.wikipedia.org/wiki/Go_(programming_language)",
		"https://en.wikipedia.org/wiki/Python",
		"https://en.wikipedia.org/wiki/Solidity",
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
	if err == nil && res.StatusCode == 200  {
		defer res.Body.Close()

		var fakeUserAgentResponse FakeUserAgentResponse
		json.NewDecoder(res.Body).Decode(&fakeUserAgentResponse)
		return fakeUserAgentResponse.Result
	}
	
	var emptySlice []string
	return emptySlice
}

func main() {
	fmt.Println("[STARTING] Executing program...")

	scrapper := colly.NewCollector()

	var articles []ArticleData

	scrapper.Limit(&colly.LimitRule{
		Delay: 1,
		RandomDelay: time.Duration(rand.Intn(10)),
	})

	userAgentList := GetUserAgentList()

	// scrapper.SetProxy("http://server:port")

	scrapper.OnRequest(func(r *colly.Request) {
		fmt.Println("[INFO] Requesting page: ", r.URL.String())
		r.Headers.Set("User-Agent", GetRandomUserAgent(userAgentList))
	})

	scrapper.OnResponse(func(r *colly.Response) {
		fmt.Println("[INFO] Got a response from page: ", r.Request.URL)
	})

	var article ArticleData
	scrapper.OnHTML("main", func(e *colly.HTMLElement) {
		article.Name = e.ChildText("span.mw-page-title-main")
		article.Name = strings.TrimSpace(article.Name)
		// TODO: improve selector to pick intro only mf-section-0
		article.Intro = e.ChildText("div.mw-parser-output")
		article.Intro = strings.TrimSpace(article.Intro)

		re := regexp.MustCompile(`\[\d+]`)
		article.Intro = re.ReplaceAllString(article.Intro, "")

		articles = append(articles, article)
	})

	scrapper.OnError(func(r * colly.Response, e error) {
		log.Fatal("An error occurred: ", e)
	})
		
	for _, page := range urls_to_scrap {
		scrapper.Visit(page)
	}

	js, err := json.MarshalIndent(articles, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[INFO] Writing data to file...")
	if err := os.WriteFile("articles-data.json", js, 0664); err == nil {
		fmt.Println("[END] Data written to file successfully!")
	}
}

