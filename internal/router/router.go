package router

import (
	"blog/config"
	"blog/internal/controller"
	"blog/internal/dao"
	"blog/internal/middleware"
	"blog/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(cfg config.Config, db *gorm.DB) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Logger(), middleware.Recovery())
	registerRoutes(engine, cfg, db, false)
	return engine
}

// NewAdmin creates the HTTP server dedicated to the administrator backend.
// It also serves the built frontend so the admin UI and its API remain same-origin.
func NewAdmin(cfg config.Config, db *gorm.DB) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Logger(), middleware.Recovery())
	registerRoutes(engine, cfg, db, true)
	return engine
}

func registerRoutes(engine *gin.Engine, cfg config.Config, db *gorm.DB, adminServer bool) {

	uploadConfigured := strings.TrimSpace(cfg.Upload.URL) != "" && strings.TrimSpace(cfg.Upload.Dir) != ""
	if uploadConfigured {
		engine.Static(cfg.Upload.URL, cfg.Upload.Dir)
	}

	healthController := controller.NewHealthController()
	userDAO := dao.NewUserDAO(db)
	categoryDAO := dao.NewCategoryDAO(db)
	articleDAO := dao.NewArticleDAO(db)
	commentDAO := dao.NewCommentDAO(db)
	tagDAO := dao.NewTagDAO(db)
	activityDAO := dao.NewActivityDAO(db)
	authService := service.NewAuthService(userDAO, cfg.JWT.Secret, cfg.JWT.Expire)
	categoryService := service.NewCategoryService(categoryDAO, articleDAO)
	articleService := service.NewArticleService(articleDAO, categoryDAO, commentDAO, tagDAO)
	commentService := service.NewCommentService(commentDAO, articleDAO)
	adminCommentService := service.NewAdminCommentService(commentDAO)
	tagService := service.NewTagService(tagDAO)
	activityService := service.NewActivityService(activityDAO)
	uploadService := service.NewUploadService(cfg.Upload.Dir, cfg.Upload.URL, cfg.Upload.MaxSizeBytes)
	userService := service.NewUserService(userDAO)
	adminService := service.NewAdminService(articleDAO, categoryDAO, userDAO, commentDAO)
	siteStatsService := service.NewSiteStatsService(articleDAO, categoryDAO, tagDAO)
	authController := controller.NewAuthController(authService)
	categoryController := controller.NewCategoryController(categoryService)
	articleController := controller.NewArticleController(articleService)
	commentController := controller.NewCommentController(commentService)
	adminCommentController := controller.NewAdminCommentController(adminCommentService)
	tagController := controller.NewTagController(tagService)
	activityController := controller.NewActivityController(activityService)
	uploadController := controller.NewUploadController(uploadService)
	userController := controller.NewUserController(userService)
	adminController := controller.NewAdminController(adminService)
	siteStatsController := controller.NewSiteStatsController(siteStatsService)

	registerHealthRoute(engine, healthController)

	api := engine.Group("/api")
	registerAuthRoutes(api, authController)
	registerPublicRoutes(api, categoryController, tagController, articleController, commentController, siteStatsController, activityController)
	if !uploadConfigured {
		uploadController = nil
	}
	registerUserRoutes(api, cfg, authController, userController, categoryController, articleController, commentController, uploadController, userDAO)
	if adminServer {
		registerAdminRoutes(api, cfg, adminController, adminCommentController, userController, categoryController, tagController, articleController, userDAO)
	}

	distDir := "web/dist"
	if _, err := os.Stat(distDir); err == nil {
		engine.Static("/assets", filepath.Join(distDir, "assets"))
		engine.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))
		engine.StaticFile("/favicon.svg", filepath.Join(distDir, "favicon.svg"))
		engine.StaticFile("/icons.svg", filepath.Join(distDir, "icons.svg"))
		engine.StaticFile("/author-profile.png", filepath.Join(distDir, "author-profile.png"))
		engine.StaticFile("/blog-background.jpg", filepath.Join(distDir, "blog-background.jpg"))
		engine.StaticFile("/blog-background-wide.jpg", filepath.Join(distDir, "blog-background-wide.jpg"))
		engine.StaticFile("/blog-background-top.png", filepath.Join(distDir, "blog-background-top.png"))
		engine.StaticFile("/shorekeeper-chibi.png", filepath.Join(distDir, "shorekeeper-chibi.png"))
		engine.StaticFile("/shorekeeper-chibi-source.png", filepath.Join(distDir, "shorekeeper-chibi-source.png"))

		engine.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/uploads") {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found", "data": nil})
				return
			}
			c.File(filepath.Join(distDir, "index.html"))
		})
	}

}
