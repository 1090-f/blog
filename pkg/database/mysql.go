package database

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"blog/config"
	"blog/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ensureDB 创建数据库（如果不存在）
func ensureDB(cfg config.DatabaseConfig) {
	// 不指定数据库名连接 MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[init] 连接 MySQL 失败: %v", err)
		return
	}

	// 创建数据库
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Name)
	if err := db.Exec(sql).Error; err != nil {
		log.Printf("[init] 创建数据库失败: %v", err)
		return
	}

	log.Printf("[init] 数据库 %s 已就绪", cfg.Name)
}

func MustOpen(cfg config.DatabaseConfig) *gorm.DB {
	// 先确保数据库存在
	ensureDB(cfg)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.ArticleTag{},
		&model.Comment{},
	); err != nil {
		panic(err)
	}

	log.Println("[init] 数据表已同步")
	if err := ensureGuestCommentSchema(db); err != nil {
		panic(err)
	}

	return db
}

// ensureGuestCommentSchema makes user_id nullable for databases created before
// anonymous comments were supported. AutoMigrate does not always relax an
// existing NOT NULL constraint, so only those legacy schemas need a DDL change.
func ensureGuestCommentSchema(db *gorm.DB) error {
	columns, err := db.Migrator().ColumnTypes(&model.Comment{})
	if err != nil {
		return err
	}

	for _, column := range columns {
		if column.Name() != "user_id" {
			continue
		}

		nullable, ok := column.Nullable()
		if !ok {
			return fmt.Errorf("could not determine whether comments.user_id is nullable")
		}
		if nullable {
			return nil
		}

		return db.Exec("ALTER TABLE comments MODIFY COLUMN user_id BIGINT UNSIGNED NULL").Error
	}

	return fmt.Errorf("comments.user_id column not found")
}

func EnsureAdmin(db *gorm.DB, cfg config.AdminBootstrapConfig) error {
	if !cfg.Enabled {
		return nil
	}

	username := strings.TrimSpace(cfg.Username)
	password := strings.TrimSpace(cfg.Password)
	nickname := strings.TrimSpace(cfg.Nickname)
	if username == "" || password == "" || nickname == "" {
		return errors.New("admin bootstrap config is incomplete")
	}

	var existing model.User
	err := db.Where("role = ?", "admin").First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Nickname: nickname,
		Role:     "admin",
		Status:   1,
	}

	return db.Create(admin).Error
}
