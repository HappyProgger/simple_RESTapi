package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

var cfg Config

type Config struct {
	App App `yaml:"app"`
	DB  DB  `yaml:"db"`
}

type App struct {
	Adres string `yaml:"adres"`
	Port  int    `yaml:"port"`
}

type DB struct {
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
	Name     string `yaml:"db_table"`
	Adres    string `yaml:"adres"`
	Port     int    `yaml:"port"`
}

func main() {
	// Чтение файла конфигурации
	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		fmt.Errorf("The config could not be read from the file")
	}

}
