package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type App struct {
	HTTPServer HTTPServer `mapstructure:"http_server"`

	LogLevel string `mapstructure:"log_level"`
}

func (ac App) Validate() error {
	if err := ac.HTTPServer.Validate(); err != nil {
		return err
	}
	if ac.LogLevel == "" {
		return errors.New("log_level is required")
	}
	return nil
}

const prefix = "ASM"

func LoadAndValidate() (*App, error) {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(prefix)

	viper.SetConfigFile(path)

	var cfg App
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, err
		}

		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
