package routes

import (
	"fmt"
	"time"

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

	r.GET("test",func(ctx *gin.Context) {
		ctx.SetCookie("test","value",60,"/","localhost",false,true)
		ctx.JSON(200,gin.H {
			"message":ctx.Request.URL,
		})
	})

	r.GET("testv",func(ctx *gin.Context) {
		c,_ := ctx.Cookie("test")
		ctx.JSON(200,gin.H {
			"message":c,
		})
	})

	r.GET("test2",func(ctx *gin.Context) {
		now := time.Now()
		tomorrow := time.Now().Add(time.Hour * 24)

		fmt.Println(now.After(tomorrow))

		ctx.JSON(200,gin.H {
			"message":"c",
		})
	})



	// all admin routes go here
	AdminRoutes(r)
	PublicRoutes(r)

}
