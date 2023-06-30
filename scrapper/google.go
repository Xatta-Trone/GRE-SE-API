package scrapper

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

// GetGoogleResultAndSave goes to google and retrieves the google result and saves to db
func GetGoogleResultAndSave(db *sqlx.DB, word model.Result) {
	time.Sleep(500 * time.Millisecond)

	utils.PrintS(fmt.Sprintf("Getting %v - %s from google \n", word.ID, word.Word))

	// get the google scrapper url

	googleUrl := os.Getenv("GOOGLE_URL")

	if googleUrl == "" {
		utils.PrintR("No google url found")
		return
	}

	// we have the google url
	res, err := http.Get(fmt.Sprintf("%s/%s", googleUrl, word.Word))

	if err != nil {
		utils.Errorf(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		// insert into the db
		body, _ := io.ReadAll(res.Body)
		_, err := db.Exec("Update wordlist set google=?,is_google_parsed=1,updated_at=now() where id = ? ", string(body), word.ID)

		if err != nil {
			utils.Errorf(err)
			return
		}

		utils.PrintG(fmt.Sprintf("Inserted %v - %s from google \n", word.ID, word.Word))

	}

	if res.StatusCode != http.StatusOK {
		_, err := db.Exec("Update wordlist set google_try= google_try+1,updated_at=now() where id = ? ", word.ID)

		if err != nil {
			utils.Errorf(err)
		}
		utils.PrintR(fmt.Sprintf("Updated Not found %v - %s from google \n", word.ID, word.Word))

	}


}

func GetGoogleResultAndSaveWithWG(db *sqlx.DB, word model.Result, wg *sync.WaitGroup) {
	// time.Sleep(500 * time.Millisecond)
	defer wg.Done()

	utils.PrintS(fmt.Sprintf("Getting %v - %s from google \n", word.ID, word.Word))

	// get the google scrapper url

	googleUrl := os.Getenv("GOOGLE_URL")

	if googleUrl == "" {
		utils.PrintR("No google url found")
		return
	}

	// we have the google url
	res, err := http.Get(fmt.Sprintf("%s/%s", googleUrl, word.Word))

	if err != nil {
		utils.Errorf(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		// insert into the db
		body, _ := io.ReadAll(res.Body)
		_, err := db.Exec("Update wordlist set google=?,is_google_parsed=1,updated_at=now() where id = ? ", string(body), word.ID)

		if err != nil {
			utils.Errorf(err)
			return
		}

		utils.PrintG(fmt.Sprintf("Inserted %v - %s from google \n", word.ID, word.Word))

	}

	if res.StatusCode != http.StatusOK {
		_, err := db.Exec("Update wordlist set google_try= google_try+1,updated_at=now() where id = ? ", word.ID)

		if err != nil {
			utils.Errorf(err)
		}
		utils.PrintR(fmt.Sprintf("Updated Not found %v - %s from google \n", word.ID, word.Word))

	}

}

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
			utils.Errorf(err)
		}

		defer res.Body.Close()

		// fmt.Println(string(body))

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			body, _ := ioutil.ReadAll(res.Body)
			_, err := database.Gdb.Exec("Update wordlist set google=?,is_google_parsed=1 where id = ? ", string(body), word.ID)

			if err != nil {
				utils.Errorf(err)
			}

			// fmt.Printf("Inserted %v - %s from google \n", word.ID, word.Word)
			str := fmt.Sprintf("Inserted %v - %s from google \n", word.ID, word.Word)
			utils.PrintG(str)
		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set google_try= google_try+1 where id = ? ", word.ID)

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
	fmt.Println("Done the google scrapper")

}
