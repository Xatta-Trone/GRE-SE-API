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

	database.Gdb = database.InitializeDB()

	defer database.Gdb.Close()

	// GetChatGpt()

	// populate the words
	// readRemoteFile()

	// var wg sync.WaitGroup
	// // // populate google result
	// wg.Add(1)
	// // go scrapper.GetGoogleResult(&wg)
	// go scrapper.GetWikiResult(&wg)
	//  // go scrapper.GetThesaurusResult(&wg)
	// // go scrapper.GetWordsResult(&wg)
	// // // go scrapper.GetNinjaResult(&wg)

	// wg.Wait()

	// ReadAndImportNamedCsv("Barrons-333.csv", "Barron's 333")

	processor.ReadTableAndProcessWord("aggregate")

	

	//
	fmt.Println("All done")
	elapsed := time.Since(start)
	log.Printf("Total time took %s", elapsed)

}
