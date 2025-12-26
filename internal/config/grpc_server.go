package config

type GRPCServer struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Reflect bool   `mapstructure:"reflect"`
}
