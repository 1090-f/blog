package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置，包含服务器、数据库、JWT、上传和管理员初始化等全部配置项。
type Config struct {
	Server         ServerConfig         `mapstructure:"server"`
	AdminServer    ServerConfig         `mapstructure:"admin_server"`
	Database       DatabaseConfig       `mapstructure:"database"`
	JWT            JWTConfig            `mapstructure:"jwt"`
	Upload         UploadConfig         `mapstructure:"upload"`
	AdminBootstrap AdminBootstrapConfig `mapstructure:"admin_bootstrap"`
}

// ServerConfig HTTP 服务端口配置。
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig MySQL 数据库连接参数。
type DatabaseConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Name string `mapstructure:"name"`
}

// JWTConfig JWT 签名密钥与令牌过期时间配置。
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

// UploadConfig 文件上传目录、访问 URL 前缀及大小限制配置。
type UploadConfig struct {
	Dir          string `mapstructure:"dir"`
	URL          string `mapstructure:"url"`
	MaxSizeBytes int64  `mapstructure:"max_size_bytes"`
}

// AdminBootstrapConfig 首次运行时自动创建的管理员账号配置。
type AdminBootstrapConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Nickname string `mapstructure:"nickname"`
}

// Load 读取配置文件、默认值和环境变量，生成应用配置。
func Load() (Config, error) {
	v := viper.New()
	v.SetDefault("server.port", "8080") // 设置默认值
	v.SetDefault("admin_server.port", "8081")
	v.SetDefault("jwt.expire", 7200)
	v.SetDefault("upload.max_size_bytes", int64(5*1024*1024))
	v.SetEnvPrefix("blog")                             // 为环境变量添加前缀:BLOG_
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // . → _ 替换
	v.AutomaticEnv()                                   // 自动读取环境变量
	bindEnv(v,
		"server.port",
		"admin_server.port",
		"database.host",
		"database.port",
		"database.user",
		"database.pass",
		"database.name",
		"jwt.secret",
		"jwt.expire",
		"upload.dir",
		"upload.url",
		"upload.max_size_bytes",
		"admin_bootstrap.enabled",
		"admin_bootstrap.username",
		"admin_bootstrap.password",
		"admin_bootstrap.nickname",
	)

	configFile := strings.TrimSpace(os.Getenv("BLOG_CONFIG_FILE"))
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")
	}

	var cfg Config
	if err := readConfig(v, configFile != ""); err != nil {
		return cfg, err
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// MustLoad 加载应用配置；加载失败时终止程序。
func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

// 按配置文件是否显式指定的规则读取配置。
func readConfig(v *viper.Viper, explicitConfigFile bool) error {
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) && !explicitConfigFile {
			return nil
		}
		return err
	}
	return nil
}

// 将配置键绑定到同名环境变量。
func bindEnv(v *viper.Viper, keys ...string) {
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}
