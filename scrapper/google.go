package scrapper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

func GetGoogleResult(wg *sync.WaitGroup) {
	defer wg.Done()

	utils.PrintS("Inside the google scrapper")

	words := []model.WordGetStruct{}
	database.Gdb.Select(&words, "SELECT id, word from wordlist where is_google_parsed=0 and google_try < 6")

	for _, word := range words {

		time.Sleep(800 * time.Millisecond)
		fmt.Printf("Getting %v - %s from google \n", word.ID, word.Word)

		// res, err := http.Get(fmt.Sprintf("https://dict.gre-sentence-equivalence.com/word/%s", word.Word))
		// res, err := http.Get(fmt.Sprintf("https://dictionary-api-v7nc.onrender.com/word/%s", word.Word))
		res, err := http.Get(fmt.Sprintf("http://localhost:8080/word/%s", word.Word))

		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()

		// fmt.Println(string(body))

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			body, _ := ioutil.ReadAll(res.Body)
			_, err := database.Gdb.Exec("Update wordlist set google=?,is_google_parsed=1 where id = ? ", string(body), word.ID)

			if err != nil {
				fmt.Println(err)
			}

			// fmt.Printf("Inserted %v - %s from google \n", word.ID, word.Word)
			str := fmt.Sprintf("Inserted %v - %s from google \n", word.ID, word.Word)
			utils.PrintG(str)
		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set google_try= google_try+1 where id = ? ", word.ID)

			if err != nil {
				fmt.Println(err)
			}

		}

		if res.StatusCode == http.StatusTooManyRequests {
			// wg.Done()
			break
		}

	}

	// fmt.Println(words)
	fmt.Println("Done the google scrapper")

}
