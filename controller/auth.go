package controller

import (
	"encoding/json"
	"net/http"

	"github.com/sistemakreditasi/backend-akreditasi/config"
	"github.com/sistemakreditasi/backend-akreditasi/helper"
	"github.com/sistemakreditasi/backend-akreditasi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Fungsi login
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Gunakan koneksi MongoDB dari config
	db := config.Mongoconn
	collection := db.Collection("users")

	var foundUser models.User
	err = collection.FindOne(r.Context(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !helper.CheckPasswordHash(user.Password, foundUser.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token (gunakan fungsi CreateJWT dari helper)
	token, err := helper.CreateJWT(foundUser.Email)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Return the JWT token
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
