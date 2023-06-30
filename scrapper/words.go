package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/imroc/req/v3"
	"github.com/jmoiron/sqlx"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

// GetWordsResultAndSave goes to thesaurus and retrieves the result and saves to db
func GetWordsResultAndSaveWithWg(db *sqlx.DB, word model.Result, wg *sync.WaitGroup) {
	defer wg.Done()
	utils.PrintS(fmt.Sprintf("Getting %v - %s from wordsApi \n", word.ID, word.Word))

	client := req.C().SetCommonHeaders(map[string]string{
		"x-rapidapi-key":  os.Getenv("WORDS_API"),
		"x-rapidapi-host": "wordsapiv1.p.rapidapi.com",
	}) // Use C() to create a client.
	res, err := client.R(). // Use R() to create a request.
				Get(fmt.Sprintf("https://wordsapiv1.p.rapidapi.com/words/%s", word.Word))

	if err != nil {
		utils.Errorf(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		// insert into the db
		_, err := db.Exec("Update wordlist set words_api=?,is_words_api_parsed=1,updated_at=now() where id = ? ", res.String(), word.ID)

		if err != nil {
			utils.Errorf(err)
			return
		}

		utils.PrintG(fmt.Sprintf("Inserted %v - %s from words api \n", word.ID, word.Word))

	}

	if res.StatusCode != http.StatusOK {
		_, err := db.Exec("Update wordlist set words_api_try= words_api_try+1,updated_at=now() where id = ? ", word.ID)

		if err != nil {
			utils.Errorf(err)
		}
		utils.PrintR(fmt.Sprintf("Updated Not found %v - %s from words api \n", word.ID, word.Word))

	}

}

// GetWordsResultAndSave goes to thesaurus and retrieves the result and saves to db
func GetWordsResultAndSave(db *sqlx.DB, word model.Result) {

	utils.PrintS(fmt.Sprintf("Getting %v - %s from wordsApi \n", word.ID, word.Word))

	client := req.C().SetCommonHeaders(map[string]string{
		"x-rapidapi-key":  os.Getenv("WORDS_API"),
		"x-rapidapi-host": "wordsapiv1.p.rapidapi.com",
	}) // Use C() to create a client.
	res, err := client.R(). // Use R() to create a request.
				Get(fmt.Sprintf("https://wordsapiv1.p.rapidapi.com/words/%s", word.Word))

	if err != nil {
		utils.Errorf(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		// insert into the db
		_, err := db.Exec("Update wordlist set words_api=?,is_words_api_parsed=1,updated_at=now() where id = ? ", res.String(), word.ID)

		if err != nil {
			utils.Errorf(err)
			return
		}

		utils.PrintG(fmt.Sprintf("Inserted %v - %s from words api \n", word.ID, word.Word))

	}

	if res.StatusCode != http.StatusOK {
		_, err := db.Exec("Update wordlist set words_api_try= words_api_try+1,updated_at=now() where id = ? ", word.ID)

		if err != nil {
			utils.Errorf(err)
		}
		utils.PrintR(fmt.Sprintf("Updated Not found %v - %s from words api \n", word.ID, word.Word))

	}

}

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
				utils.Errorf(err)
				continue
			}

			str := fmt.Sprintf("Inserted %v - %s from words api \n", word.ID, word.Word)
			utils.PrintG(str)
		}

		if res.StatusCode == http.StatusNotFound {
			_, err := database.Gdb.Exec("Update wordlist set words_api_try= words_api_try+1 where id = ? ", word.ID)

			if err != nil {
				utils.Errorf(err)
			}

		} else {
			continue
		}

	}
	fmt.Println("Done the words api scrapper")

}
