package routes

import (
	"net/http"

	"github.com/sistemakreditasi/backend-akreditasi/config"
	"github.com/sistemakreditasi/backend-akreditasi/controller"
	"github.com/sistemakreditasi/backend-akreditasi/helper"
)

// URL handles incoming HTTP requests and routes them to appropriate controllers
func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}

	if config.ErrorMongoconn != nil {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : database, "+config.ErrorMongoconn.Error())
		return
	}

	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		Home(w, r)
	case method == "POST" && path == "/upload":
		controller.UploadPDF(w, r)
	case method == "GET" && path == "/download":
		controller.DownloadPDF(w, r)
	case method == "POST" && path == "/login":
		controller.Login(w, r)
	default:
		helper.ErrorResponse(w, r, http.StatusNotFound, "Not Found", "The requested resource was not found")
	}
}

// Home function is just a simple handler for root URL
func Home(respw http.ResponseWriter, req *http.Request) {
	resp := map[string]string{
		"message": "Welcome to the PDF Service",
	}
	helper.WriteJSON(respw, http.StatusOK, resp)
}
