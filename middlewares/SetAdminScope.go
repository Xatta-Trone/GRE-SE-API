package middlewares

import (
	"github.com/gin-gonic/gin"
)

func SetAdminScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Request.Form.Add("scope","test")
		// fmt.Println(c.Request.Form.Has("scope"))
		c.Next()
	}
}
