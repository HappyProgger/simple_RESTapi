package main

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigServer struct {
	Port string `env:"PORT" env-description:"server port"`
	Host string `env:"HOST" env-description:"server host"`
}

func main() {
	var cfg ConfigServer

	// Чтение конфигурации из переменных окружения
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		fmt.Println("Ошибка при чтении конфигурации:", err)
		os.Exit(1)
	}

	// Вывод конфигурации
	fmt.Printf("Host: %s\n", cfg.Host)
	fmt.Printf("Port: %s\n", cfg.Port)
}
