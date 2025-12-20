package config

import "net"

type GRPCServer struct {
	Host `mapstructure:"host"`
	Port `mapstructure:"port"`

	ReflectEnabled bool `mapstructure:"reflect_enabled"`
}

func (s GRPCServer) Validate() error {
	if err := s.Host.Validate(); err != nil {
		return err
	}
	if err := s.Port.Validate(); err != nil {
		return err
	}
	return nil
}

func (s GRPCServer) GetAddress() string {
	return net.JoinHostPort(s.Host.String(), s.Port.String())
}
