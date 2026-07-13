package router

import (
	"blog/config"
	"blog/internal/controller"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerHealthRoute(engine *gin.Engine, healthController *controller.HealthController) {
	engine.GET("/health", healthController.Check)
}

func registerPublicRoutes(api *gin.RouterGroup, cfg config.Config, categoryController *controller.CategoryController, tagController *controller.TagController, articleController *controller.ArticleController, commentController *controller.CommentController, siteStatsController *controller.SiteStatsController, activityController *controller.ActivityController, userReader middleware.UserReader) {
	api.GET("/categories", categoryController.List)
	api.GET("/tags", tagController.List)
	api.GET("/site-stats", siteStatsController.Get)
	api.GET("/site-activity", activityController.Get)
	api.GET("/articles", articleController.List)
	api.GET("/articles/latest", articleController.Latest)
	api.GET("/articles/popular", articleController.Popular)
	api.GET("/articles/:id", articleController.Detail)
	api.GET("/articles/:id/full", articleController.FullDetail)
	api.GET("/articles/:id/comments", commentController.ListByArticle)
	api.POST("/comments", middleware.OptionalAuth(cfg.JWT.Secret, userReader), commentController.Create)
}

func registerAuthRoutes(api *gin.RouterGroup, authController *controller.AuthController) {
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/login", authController.Login)
}

func registerUserRoutes(api *gin.RouterGroup, cfg config.Config, authController *controller.AuthController, categoryController *controller.CategoryController, articleController *controller.ArticleController, commentController *controller.CommentController, uploadController *controller.UploadController, userReader middleware.UserReader) {
	userGroup := api.Group("/user")
	userGroup.Use(middleware.Auth(cfg.JWT.Secret, userReader))
	userGroup.GET("/session", authController.Session)

	authenticatedGroup := api.Group("")
	authenticatedGroup.Use(middleware.Auth(cfg.JWT.Secret, userReader))
	authenticatedGroup.DELETE("/comments/:id", commentController.DeleteMine)
	authenticatedGroup.POST("/articles", articleController.Create)
	authenticatedGroup.POST("/categories", categoryController.Create)
	if uploadController != nil {
		authenticatedGroup.POST("/upload", uploadController.Upload)
	}
}

func registerAdminRoutes(api *gin.RouterGroup, cfg config.Config, adminController *controller.AdminController, adminCommentController *controller.AdminCommentController, userController *controller.UserController, categoryController *controller.CategoryController, tagController *controller.TagController, articleController *controller.ArticleController, userReader middleware.UserReader) {
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.Auth(cfg.JWT.Secret, userReader), middleware.Admin())
	adminGroup.GET("/dashboard", adminController.Dashboard)
	adminGroup.GET("/comments", adminCommentController.List)
	adminGroup.PUT("/comments/:id/status", adminCommentController.UpdateStatus)
	adminGroup.DELETE("/comments/:id", adminCommentController.Delete)
	adminGroup.GET("/users", userController.List)
	adminGroup.PUT("/users/:id/status", userController.UpdateStatus)
	adminGroup.GET("/articles", articleController.AdminList)
	adminGroup.POST("/articles", articleController.Create)
	adminGroup.PUT("/articles/:id", articleController.Update)
	adminGroup.DELETE("/articles/:id", articleController.Delete)
	adminGroup.POST("/categories", categoryController.Create)
	adminGroup.PUT("/categories/:id", categoryController.Update)
	adminGroup.DELETE("/categories/:id", categoryController.Delete)
	adminGroup.POST("/tags", tagController.Create)
	adminGroup.PUT("/tags/:id", tagController.Update)
	adminGroup.DELETE("/tags/:id", tagController.Delete)
}
