package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	var wg sync.WaitGroup
	errorChannel := make(chan error, 10)

	pages := []string{
		"https://www.linkedin.com/company/amazon/",
		"https://www.linkedin.com/company/meta/",
		"https://www.linkedin.com/company/google/",
		"https://www.linkedin.com/company/bloomberg/",
	}

	wg.Add(len(pages))

	for _, page := range pages {
		go scrape(page, &wg, errorChannel, 3)
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for err := range errorChannel {
		if err != nil {
			log.Printf("scraping error: %v", err)
		}
	}
}

func scrape(url string, wg *sync.WaitGroup, errChan chan error, retries int) {
	defer wg.Done()

	var err error
	for i := 0; i <= retries; i++ {
		if i > 0 {
			log.Printf("retrying... attempts left: %d", retries-i)
		}

		err = performScrape(url)
		if err == nil {
			return
		}

		log.Printf("failed to scrape page: %v", err)
	}

	errChan <- err
}

func performScrape(url string) error {
	log.Printf("starting to scrape %s;", url)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserAgent(getRandomUserAgent()),
	)

	allocatorContext, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	cOpts := []chromedp.ContextOption{}
	//cOpts = append(cOpts, chromedp.WithDebugf(log.Printf))

	ctx, cancel := chromedp.NewContext(allocatorContext, cOpts...)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var body string

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.OuterHTML(`body`, &body, chromedp.ByQuery),
	}

	err := chromedp.Run(ctx, tasks...)
	if err != nil {
		return err
	}

	return companyInfo(body)
}

func getRandomUserAgent() string {
	var userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/601.7.7 (KHTML, like Gecko) Version/9.1.2 Safari/601.7.7",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(userAgents))

	return userAgents[randomIndex]
}

func companyInfo(htmlContent string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return err
	}

	type CompanyInfo struct {
		About string `json:"about,omitempty"`
	}

	var matches []string
	doc.Find("p").Each(func(index int, element *goquery.Selection) {
		dataTestID, exists := element.Attr("data-test-id")
		if exists && dataTestID == "about-us__description" {
			matches = append(matches, element.Text())
		}
	})

	if len(matches) > 1 || len(matches) == 0 {
		return errors.New("multiple or no elements found, be more precise. Aborting")
	}

	info := CompanyInfo{
		About: matches[0],
	}

	jsonData, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}
