package scrapper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/xatta-trone/words-combinator/utils"
)

func ScrapVocabulary(url string) ([]string, string, error) {
	var words []string
	var titleMeta string
	var err error

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			// fmt.Println(string(r.Body))

			fmt.Println("response status code")
			fmt.Println(r.StatusCode)

			if r.StatusCode != http.StatusOK {
				fmt.Println("There was an error, ", r.Status)
				err = fmt.Errorf("%s", r.Status)
			}

			// get the title

			title := r.HTMLDoc.Find("h1.title").Text()
			title = strings.TrimSpace(title)

			fmt.Print(title)

			if len(title) > 0 {
				titleMeta = title
			}

			// get the words
			root := r.HTMLDoc.Find("ol.wordlist")

			fmt.Println(root.Length())

			if root.Length() == 1 {

				root.Children().Each(func(i int, s *goquery.Selection) {
					wordCheck := s.AttrOr("word", "")

					fmt.Println(wordCheck)

					if wordCheck != "" {
						word := strings.TrimSpace(strings.ReplaceAll(wordCheck, "\n", " "))
						words = append(words, word)

					}

				})

			}

		},
	}).Start()

	return words, titleMeta, err
}

func ScrapVocabularyOld(url string) ([]string, string, error) {
	var words []string
	var fileName string

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("ScrapVocabulary error in building request")
		return words, fileName, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0")

	resp, err := client.Do(req)
	utils.PrintG("visiting " + url)

	if err != nil {
		fmt.Println("ScrapVocabulary error in making request ")
		return words, fileName, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("ScrapVocabulary error in response status", err)
		return words, fileName, errors.New(resp.Status)

	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("ScrapVocabulary error in parsing body ", err)
		return words, fileName, err
	}

	// find the title,
	titleData := doc.Find("h1.title").Text()
	title := strings.TrimSpace(titleData)
	if len(title) > 0 {
		fileName = title
	}

	// find the words

	doc.Find("ol.wordlist").Children().Each(func(i int, s *goquery.Selection) {
		// we are inside each list element
		// <li class="entry learnable" id="entry1" word="estranged" freq="2906.44" lang="en">
		// <a class="word" href="/dictionary/estranged" title="caused to be unloved"><span class="count"></span> estranged</a>
		// <div class="definition" title="This word is learnable">caused to be unloved</div>
		// </li>

		// check if word is not null or exists
		wordCheck := s.AttrOr("word", "")
		// fmt.Println(s.AttrOr("word", ""))
		if len(wordCheck) > 0 {
			sanitizedWord := strings.TrimSpace(strings.ReplaceAll(wordCheck, "\n", " "))
			words = append(words, sanitizedWord)
		}

	})

	return words, fileName, nil
}
