package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var JWTSecret = []byte("mysecretkey") // Ganti dengan secret key yang aman

// Fungsi untuk menginisialisasi Google Drive Client
func GetDriveService() (*drive.Service, error) {
	ctx := context.Background()
	service, err := drive.NewService(ctx, option.WithCredentialsFile("path/to/credentials.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return service, err
}

// Fungsi untuk menghasilkan JWT Token
func GenerateJWT(id, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  jwt.TimeFunc().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})
	return token.SignedString(JWTSecret)
}

func GetEnv(envName string) string {
	// envFile, _ := godotenv.Read("../.env")
	// return envFile[envName]
	return os.Getenv(envName)
}
