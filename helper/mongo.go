package helper

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBInfo struct {
	DBString string
	DBName   string
}

func MongoConnect(mconn DBInfo) (db *mongo.Database, err error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mconn.DBString))
	if err != nil {
		return nil, err
	}
	return client.Database(mconn.DBName), nil
}

func CheckPasswordHash(password, hash string) bool {
	// Tambahkan fungsi untuk mengecek password dengan bcrypt (jika di-hash)
	// Gunakan bcrypt.CompareHashAndPassword
	return password == hash // Simplified; ini harus diubah untuk hash verification
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Fungsi untuk membuat JWT token
func CreateJWT(email string) (string, error) {
	// Membuat klaim JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token berlaku selama 72 jam
	})

	// Tandatangani token menggunakan kunci rahasia
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
