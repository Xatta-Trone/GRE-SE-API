package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/controllers"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/middlewares"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/services"
)

func AdminRoutes(r *gin.Engine) *gin.Engine {
	// init services
	wordRepo := repository.NewWordRepository(database.Gdb)
	wordController := controllers.NewWordController(wordRepo)

	admin := r.Group("/admin")

	admin.GET("/login", func(c *gin.Context) {
		token, _ := services.GenerateToken("dummy")

		c.JSON(200, gin.H{
			"message": token,
		})
	})

	auth := admin.Use(middlewares.AuthMiddleware())

	auth.GET("words", wordController.WordIndex)
	auth.GET("words/:id", wordController.WordById)
	auth.DELETE("words/:id", wordController.DeleteById)

	return r

}
