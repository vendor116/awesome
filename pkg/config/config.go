package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Validatable interface {
	Validate() error
}

func LoadAndValidate[T Validatable](path, prefix string) (*T, error) {
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

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
