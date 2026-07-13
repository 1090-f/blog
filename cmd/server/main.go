package server

import (
	"fmt"
	"log"
	"strings"

	"blog/config"
	"blog/internal/router"
	"blog/pkg/database"
)

func Main() {
	cfg := config.MustLoad()
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

	db := database.MustOpen(cfg.Database)
	if err := database.EnsureAdmin(db, cfg.AdminBootstrap); err != nil {
		log.Fatal(err)
	}

	publicEngine := router.New(cfg, db)
	adminEngine := router.NewAdmin(cfg, db)
	errCh := make(chan error, 2)
	go func() { errCh <- publicEngine.Run(":" + cfg.Server.Port) }()
	go func() { errCh <- adminEngine.Run(":" + cfg.AdminServer.Port) }()
	if err := <-errCh; err != nil {
		log.Fatal(fmt.Errorf("http server stopped: %w", err))
	}
}
