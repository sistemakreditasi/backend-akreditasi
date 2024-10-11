package config

import (
	"log"

	"github.com/sistemakreditasi/backend-akreditasi/helper"
)

// Ambil MongoDB connection string dari environment variable
var MongoString string = GetEnv("MONGOSTRING")

// DB info untuk MongoDB
var mongoinfo = helper.DBInfo{
	DBString: MongoString,
	DBName:   "db_akreditasi",
}

// Inisialisasi koneksi MongoDB
var Mongoconn, ErrorMongoconn = helper.MongoConnect(mongoinfo)

func InitMongo() {
	if ErrorMongoconn != nil {
		log.Fatalf("Error connecting to MongoDB: %v", ErrorMongoconn)
	} else {
		log.Println("Connected to MongoDB!")
	}
}
