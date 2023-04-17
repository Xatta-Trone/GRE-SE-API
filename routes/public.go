package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-combinator/controllers/publicController"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/middlewares"
	"github.com/xatta-trone/words-combinator/repository"
	"github.com/xatta-trone/words-combinator/services"
)

func PublicRoutes(r *gin.Engine) *gin.Engine {
	// init services
	userRepo := repository.NewUserRepository(database.Gdb)
	authService := services.NewAuthService()
	authController := publicController.NewAuthController(userRepo, authService)

	// list controller
	listRepo := repository.NewListRepository(database.Gdb)
	wordRepo := repository.NewWordRepository(database.Gdb)
	listService := services.NewListProcessorService(database.Gdb)
	listController := publicController.NewListsController(listRepo, listService, wordRepo,userRepo)

	// folders 
	folderRepo := repository.NewFolderRepository(database.Gdb)
	folderController := publicController.NewFolderController(folderRepo,listRepo,userRepo)

	public := r.Group("/")

	public.POST("/register", authController.Register)
	public.POST("login", authController.Login)

	// public lists 
	public.GET("/public-lists",listController.PublicLists)
	public.GET("/public-folders",folderController.PublicFolders)

	// public.GET("@:name", func(ctx *gin.Context) {
	// 	name := ctx.Param("name")
	// 	ctx.JSON(200, gin.H{"name": name})
	// })

	authRoutes := r.Group("/").Use(middlewares.PublicAuthMiddleware())

	authRoutes.GET("/me", authController.Me)
	authRoutes.PUT("/update", authController.Update)
	authRoutes.PATCH("/update", authController.Update)

	// lists
	authRoutes.GET("/lists", listController.Index)
	authRoutes.POST("/lists", listController.Create)
	authRoutes.GET("/lists/:slug", listController.FindOne)
	authRoutes.PUT("/lists/:slug", listController.Update)
	authRoutes.PATCH("/lists/:slug", listController.Update)
	authRoutes.DELETE("/lists-word/:slug", listController.DeleteWordInList)
	authRoutes.DELETE("/lists/:slug", listController.Delete)

	// folders 
	authRoutes.GET("/folders",folderController.Index)
	authRoutes.POST("/folders",folderController.Create)
	authRoutes.GET("/folders/:id",folderController.FindOne)
	authRoutes.PUT("/folders/:id",folderController.Update)
	authRoutes.PATCH("/folders/:id",folderController.Update)
	authRoutes.DELETE("/folders/:id",folderController.Delete)
	authRoutes.POST("/folders/:id/toggle-list",folderController.ToggleList)

	// saved items 
	// authRoutes.GET("/saved-lists")

	


	return r
}
