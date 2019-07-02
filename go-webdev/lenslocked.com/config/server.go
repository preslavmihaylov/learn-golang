package config

type ServerConfig struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Pepper  string `json:"pepper"`
	HMACKey string `json:"hmac_key"`
}

func (c ServerConfig) IsProduction() bool {
	return c.Env == "prod"
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:    8080,
		Env:     "dev",
		Pepper:  "some-pepper",
		HMACKey: "secret-hmac-key",
	}
}
