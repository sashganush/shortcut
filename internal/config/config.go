package config

type Config struct {
	ServerSchema string `env:"SERVER_SCHEMA" envDefault:"http://"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1"`
	ServerPort string `env:"SERVER_PORT" envDefault:":8080"`
	BaseUrl string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	BaseUri string `env:"BASE_URI" envDefault:"/"`
}

var Cfg Config

