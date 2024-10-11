package route

import (
	"github.com/gorilla/mux"
	"github.com/sistemakreditasi/backend-akreditasi/controller"
)

func RegisterRoutes(r *mux.Router) {
	// File upload/download routes
	r.HandleFunc("/upload/pdf", controller.UploadPDF).Methods("POST")
	r.HandleFunc("/download/pdf", controller.DownloadPDF).Methods("GET")
}
