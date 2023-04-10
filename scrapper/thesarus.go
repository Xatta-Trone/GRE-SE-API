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

// GetThesaurusResultAndSave goes to thesaurus and retrieves the result and saves to db
func GetThesaurusResultAndSave(db *sqlx.DB, word model.Result) {

	utils.PrintS(fmt.Sprintf("Getting %v - %s from thesaurus \n", word.ID, word.Word))

	// get the scrapper url

	url := os.Getenv("THESAURUS_URL")

	if url == "" {
		utils.PrintR("No thesaurus url found")
		return
	}

	// we have the thesaurus url
	res, err := http.Get(fmt.Sprintf("%s/%s", url, word.Word))

	if err != nil {
		utils.Errorf(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		// insert into the db
		body, _ := io.ReadAll(res.Body)
		_, err := db.Exec("Update wordlist set thesaurus=?,is_parsed_th=1,updated_at=now() where id = ? ", string(body), word.ID)

		if err != nil {
			utils.Errorf(err)
			return
		}

		utils.PrintG(fmt.Sprintf("Inserted %v - %s from thesaurus \n", word.ID, word.Word))

	}

	if res.StatusCode == http.StatusNotFound {
		_, err := db.Exec("Update wordlist set th_try= th_try+1,updated_at=now() where id = ? ", word.ID)

		if err != nil {
			utils.Errorf(err)
		}
		utils.PrintR(fmt.Sprintf("Updated Not found %v - %s from thesaurus \n", word.ID, word.Word))

	}

	if res.StatusCode == http.StatusTooManyRequests {
		color.Red("Too many attempts :: google")
		time.Sleep(4 * time.Minute)
		GetThesaurusResultAndSave(db, word)
	}

}



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
			utils.Errorf(err)
		}

		defer res.Body.Close()

		// fmt.Println(string(body))

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			body, _ := ioutil.ReadAll(res.Body)
			_, err := database.Gdb.Exec("Update wordlist set thesaurus=?,is_parsed_th=1 where id = ? ", string(body), word.ID)

			if err != nil {
				utils.Errorf(err)
			}

			utils.PrintG(fmt.Sprintf("Inserted %v - %s from thesaurus \n", word.ID, word.Word))
		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set th_try= th_try+1 where id = ? ", word.ID)

			if err != nil {
				utils.Errorf(err)
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
