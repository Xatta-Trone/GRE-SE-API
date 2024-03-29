package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/gookit/validate"
	"github.com/joho/godotenv"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/processor"
	"github.com/xatta-trone/words-combinator/routes"
	"github.com/xatta-trone/words-combinator/utils"
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

	// ====================================
	// ENV
	// ====================================

	// const projectDirName = "words-combinator"
	// re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	// cwd, _ := os.Getwd()
	// rootPath := re.Find([]byte(cwd))
	err := godotenv.Load(".env")

	// err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")

	}

	// ====================================
	// DB
	// ====================================
	database.Gdb = database.InitializeDB()

	defer database.Gdb.Close()

	// init seeder
	database.InitSeeder(database.Gdb)

	// init services
	// services.NewWordService(database.Gdb)

	// ====================================
	// Process words data
	// ====================================

	if database.Gdb != nil {
		go processor.UpdateUnUpdatedWords(database.Gdb)
	}

	// ====================================
	// CRON
	// ====================================

	cron := gocron.NewScheduler(time.UTC)

	cron.Every(1).Hours().Do(func() {
		fmt.Println(time.Now())
		utils.PrintG("Processing started")
		go processor.ProcessPendingWords(database.Gdb)
	})

	cron.StartAsync()

	// ====================================
	// ROUTING
	// ====================================

	// get release env
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// http
	gin.ForceConsoleColor()
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowCredentials = true
	// config.AllowOrigins = []string{"http://localhost:5173","*.gre-sentence-equivalence.com"}
	// config.AllowAllOrigins = true
	config.AllowOriginFunc = func(origin string) bool {
		// get allowed origin domains from env
		originsFromEnv := os.Getenv("ALLOW_ORIGIN_DOMAINS")
		origins := []string{"localhost", "gre-sentence-equivalence.com", "127.0.0.1"}
		isAllowedThisOrigin := false

		origins = append(origins, strings.Split(originsFromEnv, ",")...)
		// fmt.Println(origins,origin)

		for _, allowedOrigin := range origins {
			// fmt.Println(strings.Contains(origin, allowedOrigin))
			if strings.Contains(origin, allowedOrigin) {
				isAllowedThisOrigin = true
				break
			}
		}

		return isAllowedThisOrigin
	}

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
