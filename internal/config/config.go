package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Rest     RESTConfig
	Grpc     GrpcConfig
	Pprof    PprofConfig
	LogLevel string `mapstructure:"log_level"`
}

type RESTConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
}

type GrpcConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Reflect bool   `mapstructure:"reflect"`
}

type PprofConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	Enabled           bool          `mapstructure:"enabled"`
}

func Load() (*Config, error) {
	var path, prefix string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.StringVar(&prefix, "prefix", "", "environment variable prefix")
	flag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigFile(path)

	setDefaults()

	cfg := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, err
		}

		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}

const (
	defaultRestHost = "127.0.0.1"
	defaultRestPort = 8080

	defaultGrpcHost = "127.0.0.1"
	defaultGrpcPort = 9090

	defaultPprofHost = "127.0.0.1"
	defaultPprofPort = 6060

	defaultLogLevel = "info"

	defaultReadHeaderTimeout = time.Second * 3
	defaultShutdownTimeout   = time.Second * 3
)

func setDefaults() {
	viper.SetDefault("rest.host", defaultRestHost)
	viper.SetDefault("rest.port", defaultRestPort)
	viper.SetDefault("rest.read_header_timeout", defaultReadHeaderTimeout)
	viper.SetDefault("rest.shutdown_timeout", defaultShutdownTimeout)

	viper.SetDefault("grpc.host", defaultGrpcHost)
	viper.SetDefault("grpc.port", defaultGrpcPort)

	viper.SetDefault("pprof.host", defaultPprofHost)
	viper.SetDefault("pprof.port", defaultPprofPort)
	viper.SetDefault("pprof.read_header_timeout", defaultReadHeaderTimeout)
	viper.SetDefault("pprof.shutdown_timeout", defaultShutdownTimeout)

	viper.SetDefault("log_level", defaultLogLevel)
}
