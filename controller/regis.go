package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/sistemakreditasi/backend-akreditasi/helper"
	"github.com/sistemakreditasi/backend-akreditasi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// Register function
func Register(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var user model.User

	// Decode request body ke struct user
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	// Validasi input
	if user.Username == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	// Validasi role (pilihan role: kaprodi, dosen, staff)
	validRoles := map[string]bool{
		"kaprodi": true,
		"dosen":   true,
		"staff":   true,
	}

	if !validRoles[user.Role] {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "role tidak valid")
		return
	}

	// Validasi format email
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	// Cek apakah email sudah terdaftar
	userExists, _ := helper.GetUserFromEmail(user.Email, db)
	if user.Email == userExists.Email {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email sudah terdaftar")
		return
	}

	// Validasi password
	if strings.Contains(user.Password, " ") {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password tidak boleh mengandung spasi")
		return
	}
	if len(user.Password) < 8 {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password minimal 8 karakter")
		return
	}

	// Membuat salt dan hash password dengan Argon2
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
		return
	}
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	// Menyiapkan data pengguna yang akan dimasukkan ke database
	userData := bson.M{
		"username": user.Username,
		"email":    user.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt":     hex.EncodeToString(salt),
		"role":     user.Role,
	}

	// Menyimpan pengguna baru ke MongoDB
	insertedID, err := helper.InsertOneDoc(db, col, userData)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : insert data, "+err.Error())
		return
	}

	// Response sukses
	resp := map[string]any{
		"message":    "berhasil mendaftar",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
			"role":  user.Role,
		},
	}
	helper.WriteJSON(respw, http.StatusCreated, resp)
}
