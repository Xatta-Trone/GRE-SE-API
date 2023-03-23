package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xatta-trone/words-combinator/database"
	imp "github.com/xatta-trone/words-combinator/importer"
)

func main() {
	start := time.Now()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")

	}

	database.Gdb = database.InitializeDB()

	defer database.Gdb.Close()

	// http

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	PORT := os.Getenv("PORT")
	URL := ""

	if runtime.GOOS == "windows" {
		URL = "localhost:" + PORT
	} else {
		URL = ":" + PORT
	}

	r.Run(URL) // listen and serve on 0.0.0.0:8080

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

	imp.ReadAndImportNamedCsv("Barrons-333.csv", "Barron's 333")

	// processor.ReadTableAndProcessWord("abase")

	//
	fmt.Println("All done")
	elapsed := time.Since(start)
	log.Printf("Total time took %s", elapsed)

}
