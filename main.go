package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/o1egl/paseto"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/middlewares"
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
		token := generateToken()

		c.JSON(200, gin.H{
			"message": token,
		})
	})

	r.GET("/tokend", func(c *gin.Context) {
		token := c.Query("token")

		data := DecryptData(token)

		exp := data.Expiration.Before(time.Now())

		c.JSON(200, gin.H{
			"data":    data,
			"expired": exp,
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

	// imp.ReadAndImportNamedCsv("Barrons-333.csv", "Barron's 333")

	// processor.ReadTableAndProcessWord("abase")

	//
	fmt.Println("All done")
	elapsed := time.Since(start)
	log.Printf("Total time took %s", elapsed)

}

func generateToken() string {
	symmetricKey := []byte("LwYgz6qpagfKaEii2x3Fgb7rU7TnLBKa") // Must be 32 bytes
	now := time.Now()
	exp := now.Add(24 * time.Minute)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   "test",
		Issuer:     "test_service",
		Jti:        "123",
		Subject:    "test_subject",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	// Add custom claim    to the token
	jsonToken.Set("data", "this is a signed message")
	footer := "some footer"

	// Encrypt data
	// token, err := paseto.Encrypt(symmetricKey, jsonToken, footer)
	token, err := paseto.NewV2().Encrypt(symmetricKey, jsonToken, footer)

	if err != nil {
		panic(err)
	}

	// b, _ := hex.DecodeString("b4cbfb43df4ce210727d953e4a713307fa19bb7d9f85041438d9e11b942a37741eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")

	// fmt.Println(b)
	// privateKey := ed25519.PrivateKey(b)

	// b, _ = hex.DecodeString("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	// publicKey := ed25519.PublicKey(b)

	// fmt.Println(string(publicKey))

	// // or create a new keypair
	// // publicKey, privateKey, err := ed25519.GenerateKey(nil)

	// jsonToken := paseto.JSONToken{
	// 	Expiration: time.Now().Add(24 * time.Second),
	// }

	// // Add custom claim    to the token
	// jsonToken.Set("data", "this is a signed message")
	// footer := "some footer"

	// // Sign data
	// token, err := paseto.NewV2().Sign(privateKey, jsonToken, footer)

	if err != nil {
		panic(err)
	}

	return token
}

func DecryptData(s string) paseto.JSONToken {
	symmetricKey := []byte("LwYgz6qpagfKaEii2x3Fgb7rU7TnLBKa") // Must be 32 bytes
	// Decrypt data
	var newJsonToken paseto.JSONToken
	var newFooter string
	err := paseto.NewV2().Decrypt(s, symmetricKey, &newJsonToken, &newFooter)

	// b, _ := hex.DecodeString("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	// publicKey := ed25519.PublicKey(b)

	// // Verify data
	// var newJsonToken paseto.JSONToken
	// var newFooter string
	// err := paseto.NewV2().Verify(s, publicKey, &newJsonToken, &newFooter)

	if err != nil {
		panic(err)
	}

	// fmt.Println(newJsonToken)
	// fmt.Println(newFooter)

	return newJsonToken

}
