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

func GetThesaurusResult(wg *sync.WaitGroup) {
	defer wg.Done()

	utils.PrintS("Inside the thesaurus scrapper")

	words := []model.WordGetStruct{}
	database.Gdb.Select(&words, "SELECT id, word from wordlist where is_parsed_th=0 and th_try < 6")

	for _, word := range words {

		time.Sleep(300 * time.Millisecond)
		fmt.Printf("Getting %v - %s from thesaurus \n", word.ID, word.Word)

		res, err := http.Get(fmt.Sprintf("http://localhost:8081/w/%s", word.Word))

		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()

		// fmt.Println(string(body))

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			body, _ := ioutil.ReadAll(res.Body)
			_, err := database.Gdb.Exec("Update wordlist set thesaurus=?,is_parsed_th=1 where id = ? ", string(body), word.ID)

			if err != nil {
				fmt.Println(err)
			}

			utils.PrintG(fmt.Sprintf("Inserted %v - %s from thesaurus \n", word.ID, word.Word))
		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set th_try= th_try+1 where id = ? ", word.ID)

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
	fmt.Println("Done the thesaurus scrapper")

}
