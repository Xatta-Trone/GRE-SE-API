package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/controllers"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/middlewares"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/services"
)

type AdminPass struct {
	Password string `form:"password"`
}

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

	// folders
	folderRepo := repository.NewFolderRepository(database.Gdb)
	listRepo := repository.NewListRepository(database.Gdb)
	folderController := controllers.NewAdminFolderController(folderRepo, listRepo, usersRepo)

	// list controller
	listService := services.NewListProcessorService(database.Gdb)
	listController := controllers.NewListsController(listRepo, listService, wordRepo, usersRepo)

	admin := r.Group("/admin")

	admin.Use(middlewares.SetAdminScope())

	admin.POST("/login", func(c *gin.Context) {

		// get the key
		password := os.Getenv("ADMIN_PASS")

		if password == "" {
			panic("ADMIN_PASS not found")
		}

		var admin AdminPass

		if c.ShouldBind(&admin) == nil {
			if admin.Password == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"errors": "Password required",
				})
				return
			}

			if admin.Password != password {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"errors": "Password invalid.",
				})
				return

			} else {
				token, _ := services.GenerateToken("admin")
				c.JSON(http.StatusOK, gin.H{
					"token": token,
				})
				return

			}

		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": "something went wrong.",
		})
	})

	auth := admin.Use(middlewares.AuthMiddleware())

	// admin profile
	auth.GET("/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": map[string]any{
				"id":    1,
				"name":  "Monzurul Islam",
				"email": "monzurul.ce.buet@gmail.com",
			},
		})
	})

	auth.POST("/logout", func(c *gin.Context) {
		c.JSON(http.StatusNoContent, gin.H{
			"data": map[string]any{
				"id":    1,
				"name":  "Monzurul Islam",
				"email": "monzurul.ce.buet@gmail.com",
			},
		})
	})

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

	// folders
	admin.GET("/folders", folderController.Index)
	admin.POST("/folders", folderController.Create)
	admin.GET("/folders/:id", folderController.FindOne)
	admin.PUT("/folders/:id", folderController.Update)
	admin.PATCH("/folders/:id", folderController.Update)
	admin.DELETE("/folders/:id", folderController.Delete)
	admin.POST("/folders/:id/toggle-list", folderController.ToggleList)

	// lists
	admin.GET("/lists", listController.Index)
	admin.POST("/lists", listController.Create)
	admin.GET("/lists/:id", listController.FindOne)
	admin.PUT("/lists/:id", listController.Update)
	admin.PATCH("/lists/:id", listController.Update)
	admin.DELETE("/lists/:id", listController.Delete)
	admin.DELETE("/lists-word/:slug", listController.DeleteWordInList)

	return r

}
