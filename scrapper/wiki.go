package scrapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

func GetWikiResult(wg *sync.WaitGroup) {
	defer wg.Done()

	utils.PrintS("Inside the wiki scrapper")

	// ttl := time.Millisecond * 100

	words := []model.WordGetStruct{}
	err := database.Gdb.Select(&words, "SELECT id, word from wordlist where is_wiki_parsed=0 and wiki_try < 6")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(words))

	for _, word := range words {
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("Getting %v - %s from wiki \n", word.ID, word.Word)

		res, err := http.Get(fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en_US/%s", word.Word))

		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()

		// fmt.Println(string(body))

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			body, _ := ioutil.ReadAll(res.Body)

			var result []model.Wiki
			json.Unmarshal(body, &result)

			data,_ := json.Marshal(result[0])

			_, err := database.Gdb.Exec("Update wordlist set wiki=?,is_wiki_parsed=1 where id = ? ", string(data), word.ID)

			if err != nil {
				fmt.Println(err)
			}

			str := fmt.Sprintf("Inserted %v - %s from wiki \n", word.ID, word.Word)
			utils.PrintG(str)

			// fmt.Printf("Inserted %v - %s from wiki \n", word.ID, word.Word)

		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set wiki_try= wiki_try+1 where id = ? ", word.ID)

			if err != nil {
				fmt.Println(err)
			}

		}

		if res.StatusCode == http.StatusTooManyRequests {
			color.Red("Too many attempts :: wiki")
			time.Sleep(5 * time.Minute)
			// wg.Done()
			
			continue
		}

		// continue
		// if res.StatusCode == http.StatusTooManyRequests {
		// 	time.Sleep(ttl)
		// 	ttl = ttl * 2
		// 	continue
		// }

	}
	fmt.Println("Done the wiki scrapper")

}
