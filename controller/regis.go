package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/sistemakreditasi/backend-akreditasi/config"
	"github.com/sistemakreditasi/backend-akreditasi/helper"
	"github.com/sistemakreditasi/backend-akreditasi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/argon2"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode body request ke struct User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helper.ErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "error parsing request body: "+err.Error())
		return
	}

	// Validasi input
	if user.Username == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		helper.ErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	// Validasi role
	validRoles := map[string]bool{
		"kaprodi": true,
		"dosen":   true,
		"staff":   true,
	}

	if !validRoles[user.Role] {
		helper.ErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "role tidak valid")
		return
	}

	// Validasi email
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		helper.ErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	// Cek apakah email sudah terdaftar
	db := config.Mongoconn
	collection := db.Collection("users")
	var existingUser models.User
	err = collection.FindOne(r.Context(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		helper.ErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "email sudah terdaftar")
		return
	}

	// Validasi password
	if strings.Contains(user.Password, " ") {
		helper.ErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "password tidak boleh mengandung spasi")
		return
	}
	if len(user.Password) < 8 {
		helper.ErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "password minimal 8 karakter")
		return
	}

	// Hash password menggunakan Argon2
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: salt")
		return
	}
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	// Membuat user baru dengan ObjectID
	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Username: user.Username,
		Email:    user.Email,
		Password: hex.EncodeToString(hashedPassword),
		Role:     user.Role, // Role ditentukan oleh input dari pengguna
	}

	// Simpan user baru ke MongoDB
	_, err = collection.InsertOne(r.Context(), newUser)
	if err != nil {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: insert data, "+err.Error())
		return
	}

	// Response sukses
	resp := map[string]any{
		"message": "berhasil mendaftar",
		"data": map[string]string{
			"email": newUser.Email,
			"role":  newUser.Role,
		},
	}
	helper.WriteJSON(w, http.StatusCreated, resp)
}
