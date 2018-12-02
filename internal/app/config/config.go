package config

import "github.com/spf13/viper"

const (
	DomainSuffix   = "DOMAIN_SUFFIX"
	ServerHost     = "SERVER_HOST"
	ServerPort     = "SERVER_PORT"
	ServerProtocol = "SERVER_PROTOCOL"
)

type Config struct {
	DomainSuffix   string `json:"domain_suffix"`
	ServerHost     string `json:"server_host"`
	ServerPort     int    `json:"server_port"`
	ServerProtocol string `json:"server_protocol"`
}

func GetFromEnv() *Config {
	v := viper.New()
	v.SetDefault(DomainSuffix, ".docker.")
	v.SetDefault(ServerHost, "localhost")
	v.SetDefault(ServerPort, 5300)
	v.SetDefault(ServerProtocol, "udp")

	v.AutomaticEnv()

	return &Config{
		DomainSuffix:   v.GetString(DomainSuffix),
		ServerHost:     v.GetString(ServerHost),
		ServerPort:     v.GetInt(ServerPort),
		ServerProtocol: v.GetString(ServerProtocol),
	}
}
