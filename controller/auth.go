package controller

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/sistemakreditasi/backend-akreditasi/helper"
	"github.com/sistemakreditasi/backend-akreditasi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// Login handles user login
func Login(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privatekey string) {
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body: "+err.Error())
		return
	}

	// Validasi input
	if user.Email == "" || user.Password == "" {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	// Validasi format email
	if err = checkmail.ValidateFormat(user.Email); err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	// Cek apakah email ada di database
	var foundUser model.User
	err = db.Collection("users").FindOne(req.Context(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err == mongo.ErrNoDocuments {
		helper.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "email atau password salah")
		return
	} else if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: get email "+err.Error())
		return
	}

	// Decode salt dari string
	salt, err := hex.DecodeString(foundUser.Salt)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: salt")
		return
	}

	// Hash password menggunakan salt
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != foundUser.Password {
		helper.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "password salah")
		return
	}

	// Jika menggunakan token, generate token di sini
	tokenstring, err := helper.Encode(foundUser.ID, foundUser.Email, privatekey)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: token")
		return
	}

	// Berhasil login, kirim response
	resp := map[string]interface{}{
		"status":  "success",
		"message": "login berhasil",
		"token":   tokenstring,
		"data": map[string]string{
			"email": foundUser.Email,
			"role":  foundUser.Role,
		},
	}
	helper.WriteJSON(respw, http.StatusOK, resp)
}
