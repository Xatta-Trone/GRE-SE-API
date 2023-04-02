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
	wordGroupService := services.NewWordGroupService(database.Gdb)
	wordGroupController := controllers.NewWordGroupController(wordGroupRepo, wordGroupService)

	// users service
	usersRepo := repository.NewUserRepository(database.Gdb)
	usersController := controllers.NewUsersController(usersRepo)

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
	auth.GET("/word-groups", wordGroupController.Index)
	auth.POST("/word-groups", wordGroupController.Import)
	auth.GET("/word-groups/:id", wordGroupController.FindOne)
	auth.DELETE("/word-groups/:id", wordGroupController.DeleteById)

	// users routes
	auth.GET("/users", usersController.Index)
	auth.POST("/users", usersController.Create)
	auth.GET("/users/:id", usersController.FindOne)
	auth.PATCH("/users/:id", usersController.Update)
	auth.PUT("/users/:id", usersController.Update)
	auth.DELETE("/users/:id", usersController.Delete)
	

	return r

}
