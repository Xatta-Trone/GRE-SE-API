package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/imroc/req/v3"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

func GetWordsResult(wg *sync.WaitGroup) {

	defer wg.Done()

	utils.PrintS("Inside the words api scrapper")

	words := []model.WordGetStruct{}
	database.Gdb.Select(&words, "SELECT id, word from wordlist where is_words_api_parsed=0 and words_api_try < 6")

	totalParseInDay := 600

	for _, word := range words {

		if totalParseInDay == 0 {
			break
		}
		fmt.Printf("Getting %v - %s from words api \n", word.ID, word.Word)

		client := req.C().SetCommonHeaders(map[string]string{
			"x-rapidapi-key":  os.Getenv("WORDS_API"),
			"x-rapidapi-host": "wordsapiv1.p.rapidapi.com",
		}) // Use C() to create a client.
		res, err := client.R(). // Use R() to create a request.
					Get(fmt.Sprintf("https://wordsapiv1.p.rapidapi.com/words/%s", word.Word))

		totalParseInDay--
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			_, err := database.Gdb.Exec("Update wordlist set words_api=?,is_words_api_parsed=1 where id = ? ", res.String(), word.ID)

			if err != nil {
				fmt.Println(err)
				continue
			}

			str := fmt.Sprintf("Inserted %v - %s from words api \n", word.ID, word.Word)
			utils.PrintG(str)
		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set words_api_try= words_api_try+1 where id = ? ", word.ID)

			if err != nil {
				fmt.Println(err)
			}

		} else {
			continue
		}

	}
	fmt.Println("Done the words api scrapper")

}
