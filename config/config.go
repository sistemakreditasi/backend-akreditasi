package config

import (
	"os"
)

func GetEnv(envName string) string {
	// envFile, _ := godotenv.Read("../.env")
	// return envFile[envName]
	return os.Getenv(envName)
}

var DatabaseName = "namaDatabase"
