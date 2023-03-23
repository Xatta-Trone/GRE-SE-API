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
	"github.com/xatta-trone/words-combinator/middlewares"
	"github.com/xatta-trone/words-combinator/services"
	"github.com/xatta-trone/words-combinator/utils"
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
	gin.ForceConsoleColor()
	r := gin.Default()

	r.Use(middlewares.DummyMiddleware())

	r.GET("/ping", func(c *gin.Context) {

		letter, _ := utils.GenerateRandomString(20)

		c.JSON(200, gin.H{
			"message": "pong" + letter,
		})
	})

	r.GET("/token", func(c *gin.Context) {
		token,_ := services.GenerateToken("dummy")

		c.JSON(200, gin.H{
			"message": token,
		})
	})

	r.Use(middlewares.AuthMiddleware()).GET("/e", func(c *gin.Context) {
		// token := c.Query("string")

		data, exists := c.Get("data")

		c.JSON(200, gin.H{
			"message": data,
			"a":       exists,
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

	// csvimport.ReadAndImportNamedCsv("Barrons-333.csv", "Barron's 333")


	// processor.ReadTableAndProcessWord("abase")

	//
	fmt.Println("All done")
	elapsed := time.Since(start)
	log.Printf("Total time took %s", elapsed)

}
