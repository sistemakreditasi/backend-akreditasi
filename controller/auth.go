package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sistemakreditasi/backend-akreditasi/utils"

	"github.com/sistemakreditasi/backend-akreditasi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Credential struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp Credential
	mconn := utils.SetConnection(MongoEnv, dbname)
	var dataadmin models.User
	err := json.NewDecoder(r.Body).Decode(&dataadmin)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
		return utils.GCFReturnStruct(resp)
	}

	// Query ke MongoDB untuk mencocokkan email dan password
	var user models.User
	err = mconn.Collection(Colname).FindOne(r.Context(), bson.M{"email": dataadmin.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments || !utils.CheckPasswordHash(dataadmin.Password, user.Password) {
		resp.Message = "Email atau password salah"
		return utils.GCFReturnStruct(resp)
	}

	// Jika email dan password cocok, buat token JWT
	tokenstring, err := utils.GenerateJWT(user.Email, user.Role, os.Getenv(Privatekey))
	if err != nil {
		resp.Message = "Gagal Encode Token : " + err.Error()
		return utils.GCFReturnStruct(resp)
	}

	// Sukses
	resp.Status = true
	resp.Message = "Selamat Datang " + user.Username
	resp.Token = tokenstring
	return utils.GCFReturnStruct(resp)
}
