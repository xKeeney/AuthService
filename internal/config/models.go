package config

import "time"

type Config struct {
	Logger   LoggerConfig   `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type LoggerConfig struct {
	LogFile  string `mapstructure:"log_file"`
	LogLevel string `mapstructure:"log_level"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`

	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`

	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`

	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`

	SSLMode string `mapstructure:"ssl_mode"`

	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}
