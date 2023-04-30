package middlewares

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		var token string

		// get the token
		token = c.Query("token")

		if token == "" {
			// get the token from header
			token = c.GetHeader("X-API-TOKEN")

		}

		if token == "" {
			// get the token from bearer
			bearer := c.GetHeader("Authorization")

			if bearer != "" {
				t := strings.Split(bearer, " ")
				token = t[1]
			}

		}

		// fmt.Println(token)

		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"errors": "token missing"})
			return
		}

		// now decode token

		// get the key
		key := os.Getenv("AUTH_KEY")

		if key == "" {
			panic("AUTH_KEY not found")
		}

		symmetricKey := []byte(key) // Must be 32 bytes
		// Decrypt data
		var newJsonToken paseto.JSONToken
		var newFooter string
		err := paseto.NewV2().Decrypt(token, symmetricKey, &newJsonToken, &newFooter)

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"errors": "token mismatch"})
			return
		}

		// fmt.Println(newJsonToken)

		tokenExpired := newJsonToken.Expiration.Before(time.Now())

		if tokenExpired {
			c.AbortWithStatusJSON(401, gin.H{"errors": "token expired. Please login again"})
			return
		}

		c.Set("data", newJsonToken.Get("data"))

		c.Next()
	}

}
