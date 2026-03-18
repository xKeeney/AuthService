package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load(path string) (*Config, error) {
	v := viper.New()

	// файл
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	// env override
	v.SetEnvPrefix("AUTH")
	v.AutomaticEnv()

	// поддержка SERVER_PORT вместо server.port
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
