package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"github.com/joho/godotenv"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/routes"
)

func init() {
	// change global opts
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = false
	})
}

func main() {
	start := time.Now()

	const projectDirName = "words-combinator"
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)

	// err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")

	}

	database.Gdb = database.InitializeDB()

	defer database.Gdb.Close()

	// init seeder
	database.InitSeeder(database.Gdb)

	// init services
	// services.NewWordService(database.Gdb)

	// http
	gin.ForceConsoleColor()
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	routes.Init(r)

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
	// // go scrapper.GetWikiResult(&wg)
	//  // go scrapper.GetThesaurusResult(&wg)
	// // go scrapper.GetWordsResult(&wg)
	// // // go scrapper.GetNinjaResult(&wg)
	// go scrapper.GetMWResult(&wg)

	// wg.Wait()

	// csvimport.ReadAndImportNamedCsv("Barrons-333.csv", "Barron's 333")

	// processor.ReadTableAndProcessWord("abase")

	//
	fmt.Println("All done")
	elapsed := time.Since(start)
	log.Printf("Total time took %s", elapsed)

}
