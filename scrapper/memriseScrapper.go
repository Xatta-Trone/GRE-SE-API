package scrapper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/xatta-trone/words-combinator/model"
)

func GetMemriseSets(url string) (model.MemriseScrapper, error) {
	var model model.MemriseScrapper
	var err error = nil

	c := colly.NewCollector(
		colly.AllowedDomains("app.memrise.com", "app.memrise.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
	)

	// Find the element with class word-list
	c.OnHTML("a.level.clearfix", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		model.Urls = append(model.Urls, e.Request.AbsoluteURL(link))

	})

	c.OnHTML("h1.course-name.sel-course-name", func(h *colly.HTMLElement) {
		title := strings.TrimSpace(h.Text)

		model.Title = title
	})

	// check error

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("There was an error, ", e)
		err = e
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(url)

	return model, err

}

func ScrapMemrise(url string) ([]string, error) {
	var words []string
	var err error = nil

	c := colly.NewCollector(
		colly.AllowedDomains("app.memrise.com", "app.memrise.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
	)

	// Before making a request print "Visiting detailCollector..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting detailCollector", r.URL.String())
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("There was an error, ", e)
		err = e
	})

	// scrap words

	c.OnHTML("body", func(h *colly.HTMLElement) {

		h.DOM.Find(".thing").Each(func(i int, s *goquery.Selection) {
			word := strings.TrimSpace(strings.ReplaceAll(s.Find(".col_a").Text(), "\n", " "))

			fmt.Println(word)

			words = append(words, word)

		})

	})

	c.Visit(url)

	return words, err
}
