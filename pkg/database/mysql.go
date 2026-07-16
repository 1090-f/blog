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

// ensureDB 连接 MySQL 并确保目标数据库存在，不存在时自动创建。
func ensureDB(cfg config.DatabaseConfig) {
	// 不指定数据库名连接 MySQL，以便执行 CREATE DATABASE 语句
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

	// 数据库不存在则创建，存在则跳过
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Name)
	if err := db.Exec(sql).Error; err != nil {
		log.Printf("[init] 创建数据库失败: %v", err)
		return
	}

	log.Printf("[init] 数据库 %s 已就绪", cfg.Name)
}

// MustOpen 数据库初始化，自动建表
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

	// 自动创建/更新所有业务表结构
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

	return db
}

// EnsureAdmin 在管理员不存在时根据配置创建初始管理员账户，已存在则跳过。
func EnsureAdmin(db *gorm.DB, cfg config.AdminBootstrapConfig) error {
	// 控制是否自动创建管理员
	if !cfg.Enabled {
		return nil
	}

	username := strings.TrimSpace(cfg.Username)
	password := strings.TrimSpace(cfg.Password)
	nickname := strings.TrimSpace(cfg.Nickname)
	// 配置项不能为空
	if username == "" || password == "" || nickname == "" {
		return errors.New("admin bootstrap config is incomplete")
	}

	// 检查是否已存在管理员账号
	var existing model.User
	err := db.Where("role = ?", "admin").First(&existing).Error
	if err == nil {
		return nil // 已存在管理员，跳过创建
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 密码加密后写入数据库
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Nickname: nickname,
		Role:     "admin",
		Status:   1, // 默认启用状态
	}

	return db.Create(admin).Error
}
