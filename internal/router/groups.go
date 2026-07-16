package router

import (
	"blog/config"
	"blog/internal/controller"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerHealthRoute 注册健康检查路由：GET /health。
func registerHealthRoute(engine *gin.Engine, healthController *controller.HealthController) {
	engine.GET("/health", healthController.Check)
}

// registerPublicRoutes 注册前台公开路由：分类、标签、站点统计、文章列表/详情及游客评论。
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

// registerPublicAuthRoutes 注册前台认证路由：注册和登录。
func registerPublicAuthRoutes(api *gin.RouterGroup, authController *controller.AuthController) {
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/login", authController.Login)
}

// registerAdminAuthRoutes 注册管理端登录路由。
func registerAdminAuthRoutes(api *gin.RouterGroup, authController *controller.AuthController) {
	api.POST("/admin/login", authController.AdminLogin)
}

// registerUserRoutes 注册已登录用户的路由：获取会话信息、删除自己的评论。
func registerUserRoutes(api *gin.RouterGroup, cfg config.Config, authController *controller.AuthController, commentController *controller.CommentController, userReader middleware.UserReader) {
	userGroup := api.Group("/user")
	userGroup.Use(middleware.Auth(cfg.JWT.Secret, userReader))
	userGroup.GET("/session", authController.Session)

	authenticatedGroup := api.Group("")
	authenticatedGroup.Use(middleware.Auth(cfg.JWT.Secret, userReader))
	authenticatedGroup.DELETE("/comments/:id", commentController.DeleteMine)
}

// registerAdminRoutes 注册管理后台路由：仪表盘、评论管理、用户管理、文章 CRUD、分类 CRUD、标签 CRUD 及文件上传。
func registerAdminRoutes(api *gin.RouterGroup, cfg config.Config, authController *controller.AuthController, adminController *controller.AdminController, adminCommentController *controller.AdminCommentController, userController *controller.UserController, categoryController *controller.CategoryController, tagController *controller.TagController, articleController *controller.ArticleController, uploadController *controller.UploadController, userReader middleware.UserReader) {
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.Auth(cfg.JWT.Secret, userReader), middleware.Admin())
	adminGroup.GET("/session", authController.Session)
	adminGroup.GET("/dashboard", adminController.Dashboard)
	adminGroup.GET("/comments", adminCommentController.List)
	adminGroup.PUT("/comments/:id/status", adminCommentController.UpdateStatus)
	adminGroup.DELETE("/comments/:id", adminCommentController.Delete)
	adminGroup.GET("/users", userController.List)
	adminGroup.POST("/users/admin", userController.CreateAdmin)
	adminGroup.PUT("/users/:id/status", userController.UpdateStatus)
	adminGroup.PUT("/users/:id/role", userController.UpdateRole)
	adminGroup.GET("/articles", articleController.AdminList)
	adminGroup.GET("/articles/:id", articleController.AdminDetail)
	adminGroup.POST("/articles", articleController.Create)
	adminGroup.PUT("/articles/:id", articleController.Update)
	adminGroup.DELETE("/articles/:id", articleController.Delete)
	adminGroup.GET("/categories", categoryController.List)
	adminGroup.POST("/categories", categoryController.Create)
	adminGroup.PUT("/categories/:id", categoryController.Update)
	adminGroup.DELETE("/categories/:id", categoryController.Delete)
	adminGroup.GET("/tags", tagController.List)
	adminGroup.POST("/tags", tagController.Create)
	adminGroup.PUT("/tags/:id", tagController.Update)
	adminGroup.DELETE("/tags/:id", tagController.Delete)
	if uploadController != nil {
		adminGroup.POST("/upload", uploadController.Upload)
	}
}
