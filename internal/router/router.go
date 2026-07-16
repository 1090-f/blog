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

// New 创建前台 HTTP 服务。
func New(cfg config.Config, db *gorm.DB) *gin.Engine {
	engine := newEngine()
	registerPublicServerRoutes(engine, cfg, db)
	return engine
}

// NewAdmin 创建管理端 HTTP 服务。
func NewAdmin(cfg config.Config, db *gorm.DB) *gin.Engine {
	engine := newEngine()
	registerAdminServerRoutes(engine, cfg, db)
	return engine
}

// newEngine 创建 Gin 引擎实例，挂载日志和 panic 恢复中间件。
func newEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Logger(), middleware.Recovery())
	return engine
}

// registerPublicServerRoutes 注册前台服务的全部路由：公开接口、认证接口、用户接口及静态资源。
func registerPublicServerRoutes(engine *gin.Engine, cfg config.Config, db *gorm.DB) {
	userDAO := dao.NewUserDAO(db)
	categoryDAO := dao.NewCategoryDAO(db)
	articleDAO := dao.NewArticleDAO(db)
	commentDAO := dao.NewCommentDAO(db)
	tagDAO := dao.NewTagDAO(db)
	activityDAO := dao.NewActivityDAO(db)

	authController := controller.NewAuthController(service.NewAuthService(userDAO, cfg.JWT.Secret, cfg.JWT.Expire))
	categoryController := controller.NewCategoryController(service.NewCategoryService(categoryDAO, articleDAO))
	articleController := controller.NewArticleController(service.NewArticleService(articleDAO, categoryDAO, commentDAO, tagDAO))
	commentController := controller.NewCommentController(service.NewCommentService(commentDAO, articleDAO))
	tagController := controller.NewTagController(service.NewTagService(tagDAO))
	activityController := controller.NewActivityController(service.NewActivityService(activityDAO))
	siteStatsController := controller.NewSiteStatsController(service.NewSiteStatsService(articleDAO, categoryDAO, tagDAO))

	registerHealthRoute(engine, controller.NewHealthController())
	api := engine.Group("/api")
	registerPublicAuthRoutes(api, authController)
	registerPublicRoutes(api, cfg, categoryController, tagController, articleController, commentController, siteStatsController, activityController, userDAO)
	registerUserRoutes(api, cfg, authController, commentController, userDAO)
	registerServerAssets(engine, cfg, false)
}

// registerAdminServerRoutes 注册管理后台服务的全部路由：管理接口、认证接口及静态资源。
func registerAdminServerRoutes(engine *gin.Engine, cfg config.Config, db *gorm.DB) {
	userDAO := dao.NewUserDAO(db)
	categoryDAO := dao.NewCategoryDAO(db)
	articleDAO := dao.NewArticleDAO(db)
	commentDAO := dao.NewCommentDAO(db)
	tagDAO := dao.NewTagDAO(db)

	authController := controller.NewAuthController(service.NewAuthService(userDAO, cfg.JWT.Secret, cfg.JWT.Expire))
	categoryController := controller.NewCategoryController(service.NewCategoryService(categoryDAO, articleDAO))
	articleController := controller.NewArticleController(service.NewArticleService(articleDAO, categoryDAO, commentDAO, tagDAO))
	adminController := controller.NewAdminController(service.NewAdminService(articleDAO, categoryDAO, userDAO, commentDAO))
	adminCommentController := controller.NewAdminCommentController(service.NewAdminCommentService(commentDAO))
	userController := controller.NewUserController(service.NewUserService(userDAO))
	tagController := controller.NewTagController(service.NewTagService(tagDAO))

	var uploadController *controller.UploadController
	if strings.TrimSpace(cfg.Upload.URL) != "" && strings.TrimSpace(cfg.Upload.Dir) != "" {
		uploadController = controller.NewUploadController(service.NewUploadService(cfg.Upload.Dir, cfg.Upload.URL, cfg.Upload.MaxSizeBytes))
	}

	registerHealthRoute(engine, controller.NewHealthController())
	api := engine.Group("/api")
	registerAdminAuthRoutes(api, authController)
	registerAdminRoutes(api, cfg, authController, adminController, adminCommentController, userController, categoryController, tagController, articleController, uploadController, userDAO)
	registerServerAssets(engine, cfg, true)
}

// registerServerAssets 托管前端 SPA 静态资源、上传文件目录及运行时配置脚本。
func registerServerAssets(engine *gin.Engine, cfg config.Config, adminServer bool) {
	// 运行时配置脚本
	appMode := "public"
	if adminServer {
		appMode = "admin"
	}
	//  判断前后台，决定加载哪套UI
	engine.GET("/runtime-config.js", func(c *gin.Context) {
		c.Header("Content-Type", "application/javascript; charset=utf-8")
		c.String(http.StatusOK, "window.__BLOG_APP_MODE__ = %q;", appMode)
	})
	// 上传文件静态目录
	if strings.TrimSpace(cfg.Upload.URL) != "" && strings.TrimSpace(cfg.Upload.Dir) != "" {
		engine.Static(cfg.Upload.URL, cfg.Upload.Dir)
	}
	// 检查前端构建产物
	distDir := "web/dist"
	if _, err := os.Stat(distDir); err != nil {
		return
	}

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
