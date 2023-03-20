package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/imroc/req/v3"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

func GetNinjaResult(wg *sync.WaitGroup) {

	defer wg.Done()

	utils.PrintS("Inside the ninja scrapper")

	words := []model.WordGetStruct{}
	database.Gdb.Select(&words, "SELECT id, word from wordlist where is_parsed_ninja=0 and ninja_try < 6")

	totalParseInDay := 10000

	for _, word := range words {
		if totalParseInDay == 0 {
			break
		}
		fmt.Printf("Getting %v - %s from ninja api \n", word.ID, word.Word)

		res, err := req.R().
			Get(fmt.Sprintf("https://api.api-ninjas.com/v1/thesaurus?word=%s", word.Word))

		totalParseInDay--
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			_, err := database.Gdb.Exec("Update wordlist set ninja=?,is_parsed_ninja=1 where id = ? ", res.String(), word.ID)

			if err != nil {
				fmt.Println(err)
				continue
			}

			utils.PrintG(fmt.Sprintf("Inserted %v - %s from ninja \n", word.ID, word.Word))

		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set ninja_try= ninja_try+1 where id = ? ", word.ID)
			fmt.Println(err)

		}

		if res.StatusCode == http.StatusTooManyRequests {
			// wg.Done()
			break
		}

	}

	fmt.Println("Done the ninja api scrapper")

}
