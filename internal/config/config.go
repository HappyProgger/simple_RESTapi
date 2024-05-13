package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Env         string `yaml:"env" env-default:"local"`
		Adres       string `yaml:"adres"`
		Port        int    `yaml:"port"`
		Timeout     string `yaml:"timeout"`
		Idletimeout string `yaml:"idle_timeout"`
	} `yaml:"app"`

	DB struct {
		DbUser     string `yaml:"db_user"`
		DbPassword string `yaml:"db_password"`
		DbName    string `yaml:"db_name"`
		Adres      string `yaml:"adres"`
		Port       int    `yaml:"port"`
	} `yaml:"db"`
}

func main() {

	// fmt.Printf("Адрес приложения: %s\n", cfg.App.Adres)
	// fmt.Printf("Порт приложения: %d\n", cfg.App.Port)
	// fmt.Printf("Таймаут приложения: %s\n", cfg.App.Timeout)
	// fmt.Printf("Idle Timeout приложения: %s\n", cfg.App.Idletimeout)

	// fmt.Printf("Пользователь базы данных: %s\n", cfg.DB.DbUser)
	// fmt.Printf("Пароль базы данных: %s\n", cfg.DB.DbPassword)
	// fmt.Printf("Таблица базы данных: %s\n", cfg.DB.DbTable)
	// fmt.Printf("Адрес базы данных: %s\n", cfg.DB.Adres)
	// fmt.Printf("Порт базы данных: %d\n", cfg.DB.Port)
}

var cfg Config

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}
	// Получаем путь до конфиг-файла из env-переменной CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// Проверяем существование конфиг-файла
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	// Читаем конфиг-файл и заполняем нашу структуру
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
