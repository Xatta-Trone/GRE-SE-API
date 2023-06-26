package scrapper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/gocolly/colly"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

func ScrapQuizlet(url string) ([]string, string, error) {
	var words []string
	var fileName string
	var err error = nil

		geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			// fmt.Println(string(r.Body))

			if r.StatusCode != http.StatusOK {
				fmt.Println("There was an error, ", r.Status)
				err = fmt.Errorf("%s", r.Status)
			}

			root := r.HTMLDoc.Find(".SetPage-setContentWrapper")
			
			if root.Length() == 1 {
				// it will go through each group
				rootSet := root.Find(".SetPageTerms-termsWrapper")
				sets := rootSet.Find(".SetPageTerms-term")
				length := sets.Length()

				fmt.Println(length)

				if length > 0 {
					sets.Each(func(i int, s *goquery.Selection) {
						word := s.Children().Find(".SetPageTerm-wordText").Text()
						words = append(words, word)
					})
				}

				// hidden sets
				hidden := root.Find("div[style=\"display:none\"]").Children()

				fmt.Println("hidden set", hidden.Length())

				hidden.Each(func(i int, s *goquery.Selection) {
					// set the word // word is in every even number element
					str := strings.TrimSpace(strings.ReplaceAll(s.Text(), "\n", " "))
					if i == 0 || i%2 == 0 {
						words = append(words, str)
					}

				})

				// find the title
				titleText := root.Find("div.SetPage-titleWrapper").Text()
				title := strings.TrimSpace(titleText)
				// title = strings.ReplaceAll(title, " ", "-")
				title = strings.ReplaceAll(title, ":", "")
				if len(title) > 0 {
					fileName = title
				}

			}

		},
	}).Start()


	return words, fileName, err
}


func ScrapQuizlet2(url string) ([]string, string, error) {
	var words []string
	var fileName string
	var err error = nil

	c := colly.NewCollector(
		colly.AllowedDomains("quizlet.com", "www.quizlet.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
	)

	c.OnHTML(".SetPage-titleWrapper", func(h *colly.HTMLElement) {
		title := strings.TrimSpace(h.Text)
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, ":", "")
		fmt.Println(title)
		if len(title) > 0 {
			fileName = title
		}
	})

	// Find the element with class SetPageTerms-term
	c.OnHTML(".SetPageTerms-termsWrapper", func(e *colly.HTMLElement) {

		// find the free words
		e.DOM.Children().Find(".SetPageTerms-term").Each(func(i int, s *goquery.Selection) {
			word := s.Children().Find(".SetPageTerm-wordText").Text()

			words = append(words, word)

			// fmt.Println(word)

		})

		// now go for remaining words

		e.DOM.Find("div[style=\"display:none\"]").Children().Each(func(i int, s *goquery.Selection) {

			// set the word // word is in every even number element
			str := strings.TrimSpace(strings.ReplaceAll(s.Text(), "\n", " "))
			if i == 0 || i%2 == 0 {
				words = append(words, str)

			}

		})

	})

	// check error
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("There was an error, ", e.Error())
		err = e
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		utils.PrintG("visiting" + r.URL.String())
	})

	// Start scraping on https://quizlet.com/130371046/gre-flash-cards/
	c.Visit(url)

	fmt.Println(words, fileName, err)

	return words, fileName, err
}

func GetQuizletUrlMaps(url string) ([]model.QuizletFolder,string, error) {
	indexes := []model.QuizletFolder{}
	var title string
	var err error = nil
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			// fmt.Println(string(r.Body))

			if r.StatusCode != http.StatusOK {
				err = errors.New(r.Status)
			}

			// get the title 
			title = r.HTMLDoc.Find(".DashboardHeaderTitle-main").Text()

			fmt.Println(title)

			root := r.HTMLDoc.Find(".FolderPageSetsList-setsFeed")

			if root.Length() == 1 {
				// it will go through each group
				sets := root.Find(".UISetCard")
				length := sets.Length()

				// fmt.Println(length)

				if length > 0 {
					sets.Each(func(i int, s *goquery.Selection) {
						// get the url
						setUrl := s.Find(".UIBaseCardHeader a").AttrOr("href", "")
						fmt.Println(setUrl)
						// temp data
						model := model.QuizletFolder{ID: length, Url: setUrl}

						indexes = append(indexes, model)
						length--

					})
				}

			}
		},
	}).Start()

	return indexes,title, err
}
