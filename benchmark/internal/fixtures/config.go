package fixtures

var defaultConfig = &Config{Value: "test-config"}

type Config struct {
	Value string
}

func NewConfig() *Config {
	return defaultConfig
}
