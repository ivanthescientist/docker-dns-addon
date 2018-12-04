package config

import "github.com/spf13/viper"

const (
	// DomainSuffix - domain suffix to add to container name to produce an FQDN, e.g. ".docker." will produce "container.docker." as domain
	DomainSuffix = "DOMAIN_SUFFIX"
	// ServerHost - hostname/IP for the DNS server
	ServerHost = "SERVER_HOST"
	// ServerPort - port for the DNS server
	ServerPort = "SERVER_PORT"
	// ServerProtocol - protocol for DNS server e.g. "tcp" or "udp"
	ServerProtocol = "SERVER_PROTOCOL"
)

// Config is a struct containing all config options for the app
type Config struct {
	DomainSuffix   string `json:"domain_suffix"`
	ServerHost     string `json:"server_host"`
	ServerPort     int    `json:"server_port"`
	ServerProtocol string `json:"server_protocol"`
}

// GetFromEnv constructs a config with default values while overriding what's possible with values from environment
func GetFromEnv() *Config {
	v := viper.New()
	v.SetDefault(DomainSuffix, ".docker.")
	v.SetDefault(ServerHost, "0.0.0.0")
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
