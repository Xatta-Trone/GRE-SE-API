package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/middlewares"
)

func Init(r *gin.Engine) {

	r.Use(middlewares.DummyMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to GRE-SE",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// all admin routes go here
	AdminRoutes(r)
	PublicRoutes(r)

}
