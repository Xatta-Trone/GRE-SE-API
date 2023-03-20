package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/processor"
)

func main() {
	start := time.Now()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")

	}

	database.InitializeDB()

	// GetChatGpt()

	// populate the words
	// readRemoteFile()

	// var wg sync.WaitGroup
	// // populate google result
	// wg.Add(2)
	// go GetGoogleResult(&wg)
	// // go GetWikiResult(&wg)
	// // go GetThesaurusResult(&wg)
	// go GetWordsResult(&wg)
	// // go GetNinjaResult(&wg)

	// wg.Wait()

	// ReadAndImportNamedCsv("Barrons-333.csv", "Barron's 333")

	processor.ReadTableAndProcessWord("abase")

	fmt.Println("All done")
	elapsed := time.Since(start)
	log.Printf("Total time took %s", elapsed)

}
