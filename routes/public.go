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
	couponRepo := repository.NewCouponRepository(database.Gdb)
	authService := services.NewAuthService()
	authController := publicController.NewAuthController(userRepo, authService,couponRepo)

	// list controller
	listRepo := repository.NewListRepository(database.Gdb)
	wordRepo := repository.NewWordRepository(database.Gdb)
	listService := services.NewListProcessorService(database.Gdb)
	listController := publicController.NewListsController(listRepo, listService, wordRepo,userRepo)

	// folders 
	folderRepo := repository.NewFolderRepository(database.Gdb)
	folderController := publicController.NewFolderController(folderRepo,listRepo,userRepo)

	// learning status 
	learningStatusRepo := repository.NewLearningStatusRepository(database.Gdb)
	learningStatusController := publicController.NewLearningStatusController(learningStatusRepo)

	public := r.Group("/")

	public.POST("/register", authController.Register)
	public.POST("login", authController.Login)

	// public lists 
	lists := r.Group("/").Use(middlewares.OptionalAuthMiddleware())
	lists.GET("/public-lists",listController.PublicLists)
	lists.GET("/public-folders",folderController.PublicFolders)

	// public.GET("@:name", func(ctx *gin.Context) {
	// 	name := ctx.Param("name")
	// 	ctx.JSON(200, gin.H{"name": name})
	// })

	authRoutes := r.Group("/").Use(middlewares.PublicAuthMiddleware())

	authRoutes.GET("/me", authController.Me)
	authRoutes.POST("/lg", authController.Logout)
	authRoutes.PUT("/update", authController.Update)
	authRoutes.PATCH("/update", authController.Update)
	authRoutes.POST("upgrade-user",authController.Upgrade)
	authRoutes.GET("upgrade-user",authController.PurchaseSuccess)

	// lists
	authRoutes.GET("/lists", listController.Index)
	authRoutes.POST("/lists", listController.Create)
	lists.GET("/lists/:id", listController.FindOne)
	authRoutes.PUT("/lists/:id", listController.Update)
	authRoutes.PATCH("/lists/:id", listController.Update)
	authRoutes.DELETE("/lists-word/:id", listController.DeleteWordInList)
	authRoutes.DELETE("/lists/:id", listController.Delete)
	authRoutes.GET("/lists/:id/folders", listController.FoldersInList)
	authRoutes.POST("/lists/:id/words", listController.FindWords)

	// folders 
	authRoutes.GET("/folders",folderController.Index)
	authRoutes.POST("/folders",folderController.Create)
	lists.GET("/folders/:id",folderController.FindOne)
	authRoutes.PUT("/folders/:id",folderController.Update)
	authRoutes.PATCH("/folders/:id",folderController.Update)
	authRoutes.DELETE("/folders/:id",folderController.Delete)
	authRoutes.POST("/folders/:id/toggle-list",folderController.ToggleList)
	authRoutes.GET("/folders/:id/lists",folderController.ListsInFolder)

	// saved items 
	authRoutes.GET("/saved-lists",listController.SavedLists)
	authRoutes.POST("/saved-lists",listController.SaveListItem)
	authRoutes.DELETE("/saved-lists/:list_id",listController.DeleteSavedList)

	authRoutes.GET("/saved-folders",folderController.SavedFolders)
	authRoutes.POST("/saved-folders/:folder_id",folderController.SaveFolder)
	authRoutes.DELETE("/saved-folders/:folder_id",folderController.DeleteSaveFolder)

	// learning status 
	authRoutes.GET("learning-status/:list_id",learningStatusController.Index)
	authRoutes.POST("learning-status",learningStatusController.Update)
	authRoutes.DELETE("learning-status/:list_id",learningStatusController.Delete)

	
	

	return r
}
