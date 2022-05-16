package config

type Config struct {
	ServerSchema    string `env:"SERVER_SCHEMA" envDefault:"http://"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8080"`
	ServerPort      string `env:"SERVER_PORT" envDefault:":8080"`
	BaseUrl         string `env:"BASE_URL" envDefault:"http://0.0.0.0:8080"`
	BaseUri         string `env:"BASE_URI" envDefault:"/"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/persistent"`
}

var Cfg Config

