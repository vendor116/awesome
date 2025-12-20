package config

import (
	"errors"
	"net"
	"strconv"
)

type HTTPServer struct {
	Host `mapstructure:"host"`
	Port `mapstructure:"port"`
}

func (s HTTPServer) Validate() error {
	if err := s.Host.Validate(); err != nil {
		return err
	}
	if err := s.Port.Validate(); err != nil {
		return err
	}
	return nil
}

func (s HTTPServer) GetAddress() string {
	return net.JoinHostPort(s.Host.String(), s.Port.String())
}

type Host string

func (h Host) Validate() error {
	if h == "" {
		return errors.New("host is required")
	}
	return nil
}

func (h Host) String() string {
	return string(h)
}

type Port int

func (p Port) Validate() error {
	if p < 0 || p > 65535 {
		return errors.New("port is invalid")
	}
	return nil
}

func (p Port) String() string {
	return strconv.Itoa(int(p))
}
