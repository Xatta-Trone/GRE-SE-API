package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/imroc/req/v3"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var gdb *sqlx.DB

func main() {
	start := time.Now()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect db
	// db, err := sql.Open(os.Getenv("DB_DRIVER"), fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
	// assign to global db
	gdb = db
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to db!")

	// populate the words
	// readRemoteFile()

	var wg sync.WaitGroup
	// populate google result
	wg.Add(2)
	go GetGoogleResult(&wg)
	go GetWikiResult(&wg)
	// go GetWordsResult(&wg)

	wg.Wait()

	fmt.Println("All done")
	elapsed := time.Since(start)
	log.Printf("Total time took %s", elapsed)

}

type WordGetStruct struct {
	ID   int64
	Word string
}

func GetGoogleResult(wg *sync.WaitGroup) {
	defer wg.Done()

	// fmt.Println("Inside the google scrapper")
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Inside the google scrapper")

	words := []WordGetStruct{}
	gdb.Select(&words, "SELECT id, word from wordlist where is_google_parsed=0 and google_try < 6")

	for _, word := range words {

		time.Sleep(300 * time.Millisecond)
		fmt.Printf("Getting %v - %s from google \n", word.ID, word.Word)

		// res, err := http.Get(fmt.Sprintf("https://dict.gre-sentence-equivalence.com/word/%s", word.Word))
		res, err := http.Get(fmt.Sprintf("https://localhost:8080/word/%s", word.Word))

		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()

		// fmt.Println(string(body))

		// check status code

		if res.StatusCode == http.StatusOK {
			// insert into the db
			body, _ := ioutil.ReadAll(res.Body)
			_, err := gdb.Exec("Update wordlist set google=?,is_google_parsed=1 where id = ? ", string(body), word.ID)

			if err != nil {
				fmt.Println(err)
			}

			// fmt.Printf("Inserted %v - %s from google \n", word.ID, word.Word)

			s := color.New(color.FgGreen).Add(color.Underline)
			s.Printf("Inserted %v - %s from google \n", word.ID, word.Word)
		}

		if res.StatusCode == http.StatusNotFound {
			_, err := gdb.Exec("Update wordlist set google_try= google_try+1 where id = ? ", word.ID)

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

func GetWikiResult(wg *sync.WaitGroup) {
	defer wg.Done()
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Inside the wiki scrapper")

	// ttl := time.Millisecond * 100

	words := []WordGetStruct{}
	gdb.Select(&words, "SELECT id, word from wordlist where is_wiki_parsed=0 and wiki_try < 6")

	for _, word := range words {
		time.Sleep(200 * time.Millisecond)
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
			_, err := gdb.Exec("Update wordlist set wiki=?,is_wiki_parsed=1 where id = ? ", string(body), word.ID)

			if err != nil {
				fmt.Println(err)
			}

			s := color.New(color.FgGreen).Add(color.Underline)
			s.Printf("Inserted %v - %s from wiki \n", word.ID, word.Word)

			// fmt.Printf("Inserted %v - %s from wiki \n", word.ID, word.Word)

		}

		if res.StatusCode == http.StatusNotFound {
			_, err := gdb.Exec("Update wordlist set wiki_try= wiki_try+1 where id = ? ", word.ID)

			if err != nil {
				fmt.Println(err)
			}

		}

		if res.StatusCode == http.StatusTooManyRequests {
			// wg.Done()
			color.Red("Too many attempts :: wiki")
			break
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

func GetWordsResult(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Inside the words api scrapper")

	words := []WordGetStruct{}
	gdb.Select(&words, "SELECT id, word from wordlist where is_words_api_parsed=0 and words_api_try < 6")

	totalParseInDay := 2000

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
			_, err := gdb.Exec("Update wordlist set words_api=?,is_words_api_parsed=1 where id = ? ", res.String(), word.ID)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("Inserted %v - %s from words api \n", word.ID, word.Word)

		}

		if res.StatusCode == http.StatusNotFound {
			_, err := gdb.Exec("Update wordlist set words_api_try= words_api_try+1 where id = ? ", word.ID)

			if err != nil {
				fmt.Println(err)
			}

		} else {
			continue
		}

	}
	fmt.Println("Done the words api scrapper")

}

func readRemoteFile() {

	// truncate the table first
	_, err := gdb.Exec("TRUNCATE TABLE wordlist")

	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Get("https://raw.githubusercontent.com/Xatta-Trone/gre-words-collection/main/word-list/combined.csv")

	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	reader := csv.NewReader(res.Body)

	totalRows := 0

	for {
		data, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		_, err = gdb.Exec("INSERT IGNORE INTO wordlist(word) values (?)", data[0])

		if err != nil {
			fmt.Println(err)
		}

		totalRows++

		// fmt.Println(data[0])
	}

	fmt.Println("Total rows inserted", totalRows)

}
