package configs

import (
	"github.com/joho/godotenv"
)

var keys map[string]string

func SetEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		panic("Error cargando el archivo .env")
	}
	keys, err = godotenv.Read()
	if err != nil {
		panic("No pudo leerse el archivo .env")
	}
}

func Get(key string) string {
	return keys[key]
}
