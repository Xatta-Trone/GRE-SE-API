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

	// word group
	wordGroupRepo := repository.NewWordGroupRepository(database.Gdb)
	wordGroupController := controllers.NewWordGroupController(wordGroupRepo)

	admin := r.Group("/admin")

	admin.GET("/login", func(c *gin.Context) {
		token, _ := services.GenerateToken("dummy")

		c.JSON(200, gin.H{
			"message": token,
		})
	})

	auth := admin.Use(middlewares.AuthMiddleware())

	// words
	auth.GET("words", wordController.WordIndex)
	auth.POST("words", wordController.WordSave)
	auth.GET("words/:id", wordController.WordById)
	auth.DELETE("words/:id", wordController.DeleteById)
	auth.PUT("words/:id", wordController.UpdateById)
	auth.PATCH("words/:id", wordController.UpdateById)

	// word group csv import
	auth.GET("/word-group", wordGroupController.Index)
	auth.POST("/word-group", wordGroupController.Import)
	auth.GET("/word-group/:id", wordGroupController.FindOne)
	auth.DELETE("/word-group/:id", wordGroupController.DeleteById)

	return r

}
