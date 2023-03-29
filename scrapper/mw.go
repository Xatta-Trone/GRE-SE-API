package scrapper

import (
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

func GetMWResult(wg *sync.WaitGroup) {
	defer wg.Done()

	utils.PrintS("Inside the MW scrapper")

	// ttl := time.Millisecond * 100

	words := []model.WordGetStruct{}
	err := database.Gdb.Select(&words, "SELECT id, word from wordlist where is_parsed_mw=0 and mw_try < 6")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(words))

	for _, word := range words {
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("Getting %v - %s from MW \n", word.ID, word.Word)

		res, err := http.Get(fmt.Sprintf("http://localhost:8081/mw/%s", word.Word))

		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()

		// fmt.Println(string(body))

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			body, _ := ioutil.ReadAll(res.Body)


			_, err := database.Gdb.Exec("Update wordlist set mw=?,is_parsed_mw=1 where id = ? ", string(body), word.ID)

			if err != nil {
				fmt.Println(err)
			}

			str := fmt.Sprintf("Inserted %v - %s from MW \n", word.ID, word.Word)
			utils.PrintG(str)

			// fmt.Printf("Inserted %v - %s from wiki \n", word.ID, word.Word)

		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set mw_try= mw_try+1 where id = ? ", word.ID)

			if err != nil {
				fmt.Println(err)
			}

		}

		if res.StatusCode == http.StatusTooManyRequests {
			color.Red("Too many attempts :: wiki")
			time.Sleep(3 * time.Minute)
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
	fmt.Println("Done the MW scrapper")

}
