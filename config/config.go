package config

type Config struct {
	DatabasePath string
	ServerPort   string
}

func LoadConfig() *Config {
	return &Config{
		DatabasePath: "whatsapp_sessions.db",
		ServerPort:   ":8080",
	}
}
