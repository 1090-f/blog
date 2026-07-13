package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig         `mapstructure:"server"`
	AdminServer    ServerConfig         `mapstructure:"admin_server"`
	Database       DatabaseConfig       `mapstructure:"database"`
	JWT            JWTConfig            `mapstructure:"jwt"`
	Upload         UploadConfig         `mapstructure:"upload"`
	AdminBootstrap AdminBootstrapConfig `mapstructure:"admin_bootstrap"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type UploadConfig struct {
	Dir          string `mapstructure:"dir"`
	URL          string `mapstructure:"url"`
	MaxSizeBytes int64  `mapstructure:"max_size_bytes"`
}

type AdminBootstrapConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Nickname string `mapstructure:"nickname"`
}

func Load() (Config, error) {
	v := viper.New()
	v.SetDefault("server.port", "8080")
	v.SetDefault("admin_server.port", "8081")
	v.SetDefault("jwt.expire", 7200)
	v.SetDefault("upload.max_size_bytes", int64(5*1024*1024))
	v.SetEnvPrefix("blog")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
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

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

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

func bindEnv(v *viper.Viper, keys ...string) {
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}
