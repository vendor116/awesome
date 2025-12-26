package config

import (
	"errors"
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
