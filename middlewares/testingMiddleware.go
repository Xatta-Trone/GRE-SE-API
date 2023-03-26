package middlewares

import (
	"github.com/gin-gonic/gin"
)

func DummyMiddleware() gin.HandlerFunc {
  
  return func(c *gin.Context) {
	// fmt.Println("loginggg")
    c.Next()
  }
}