package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env"`
	Postgres   `yaml:"postgres"`
	HTTPServer `yaml:"http_server"`
	Redis      `yaml:"redis"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:3000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Postgres struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Address  string `yaml:"address" env-required:"true"`
	Port     string `yaml:"port" env-default:"postgres"`
	Ssl      bool   `yaml:"ssl" env-default:"false"`
	Db       string `yaml:"db" env-required:"true"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

func MustLoad() Config {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath, isExist := os.LookupEnv("CONFIG_PATH")

	if configPath == "" && !isExist {
		log.Fatal("CONFIG_PATH is not set")
	}

	//check does file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cnf Config
	if err := cleanenv.ReadConfig(configPath, &cnf); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cnf
}
