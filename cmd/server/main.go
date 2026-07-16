package server

import (
	"fmt"
	"log"
	"strings"

	"blog/config"
	"blog/internal/router"
	"blog/pkg/database"
)

// Main 启动 HTTP 服务。
func Main() {
	// 读取yaml配置
	cfg := config.MustLoad()
	// 校验关键配置
	if strings.TrimSpace(cfg.JWT.Secret) == "" {
		log.Fatal("jwt secret must not be empty")
	}
	if cfg.JWT.Expire <= 0 {
		log.Fatal("jwt expire must be greater than zero")
	}
	if strings.TrimSpace(cfg.Server.Port) == "" || strings.TrimSpace(cfg.AdminServer.Port) == "" {
		log.Fatal("server ports must not be empty")
	}
	if cfg.Server.Port == cfg.AdminServer.Port {
		log.Fatal("server.port and admin_server.port must be different")
	}

	// 连接数据库
	db := database.MustOpen(cfg.Database)
	// 管理员账户是否存在
	if err := database.EnsureAdmin(db, cfg.AdminBootstrap); err != nil {
		log.Fatal(err)
	}
	// 创建前台服务
	publicEngine := router.New(cfg, db)
	// 创建后台服务
	adminEngine := router.NewAdmin(cfg, db)
	// 启动前后台服务并监听错误
	errCh := make(chan error, 2)
	go func() { errCh <- publicEngine.Run(":" + cfg.Server.Port) }()
	go func() { errCh <- adminEngine.Run(":" + cfg.AdminServer.Port) }()
	if err := <-errCh; err != nil {
		log.Fatal(fmt.Errorf("http server stopped: %w", err))
	}
}
