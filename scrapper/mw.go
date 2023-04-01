package scrapper

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

// GetMWResultAndSave goes to thesaurus and retrieves the result and saves to db
func GetMWResultAndSave(db *sqlx.DB, word model.Result) {

	utils.PrintS(fmt.Sprintf("Getting %v - %s from MW \n", word.ID, word.Word))

	// get the scrapper url

	url := os.Getenv("MW_URL")

	if url == "" {
		utils.PrintR("No mw url found")
		return
	}

	// we have the mw url
	res, err := http.Get(fmt.Sprintf("%s/%s", url, word.Word))

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		// insert into the db
		body, _ := io.ReadAll(res.Body)
		_, err := db.Exec("Update wordlist set mw=?,is_parsed_mw=1,updated_at=now() where id = ? ", string(body), word.ID)

		if err != nil {
			fmt.Println(err)
			return
		}

		utils.PrintG(fmt.Sprintf("Inserted %v - %s from MW \n", word.ID, word.Word))

	}

	if res.StatusCode == http.StatusNotFound {
		_, err := db.Exec("Update wordlist set mw_try= mw_try+1,updated_at=now() where id = ? ", word.ID)

		if err != nil {
			fmt.Println(err)
		}
		utils.PrintR(fmt.Sprintf("Updated Not found %v - %s from MW \n", word.ID, word.Word))

	}

	if res.StatusCode == http.StatusTooManyRequests {
		color.Red("Too many attempts :: MW")
		time.Sleep(4 * time.Minute)
		GetMWResultAndSave(db, word)
	}

}


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
