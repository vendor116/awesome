package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type App struct {
	RESTServer  HTTPServer `mapstructure:"rest_server"`
	GRPCServer  GRPCServer `mapstructure:"grpc_server"`
	PprofServer HTTPServer `mapstructure:"pprof_server"`

	LogLevel     string `mapstructure:"log_level"`
	PprofEnabled bool   `mapstructure:"pprof_enabled"`
}

func (a App) Validate() error {
	if err := a.RESTServer.Validate(); err != nil {
		return err
	}
	if err := a.GRPCServer.Validate(); err != nil {
		return err
	}
	if a.LogLevel == "" {
		return errors.New("log_level is required")
	}
	return nil
}

func Load[T any](path, prefix string) (*T, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigFile(path)

	var cfg T
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, err
		}

		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
