package config

import (
	"errors"
	"time"
)

type HTTPServer struct {
	Host              string        `mapstructure:"host"`
	Port              string        `mapstructure:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
}

func (hs HTTPServer) Validate() error {
	if hs.Host == "" {
		return errors.New("host is required")
	}
	if hs.Port == "" {
		return errors.New("port is required")
	}
	if hs.ReadHeaderTimeout == 0 {
		return errors.New("read_header_timeout is required")
	}
	return nil
}
