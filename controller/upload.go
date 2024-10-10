package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sistemakreditasi/backend-akreditasi/config"
	"github.com/sistemakreditasi/backend-akreditasi/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/drive/v3"
)

func UploadPDF(w http.ResponseWriter, r *http.Request) {
	// Get file from request
	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Initialize Google Drive service
	driveService := config.GetDriveService()

	// Upload file to Google Drive
	driveFile := &drive.File{
		Name:    handler.Filename,
		Parents: []string{"<FOLDER_ID>"}, // Ganti dengan ID folder Google Drive
	}

	uploadedFile, err := driveService.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		http.Error(w, "Failed to upload file to Google Drive", http.StatusInternalServerError)
		return
	}

	// File uploaded successfully, simpan ke database
	fileURL := fmt.Sprintf("https://drive.google.com/file/d/%s/view", uploadedFile.Id)

	akreditasi := models.Akreditasi{
		ID:               primitive.NewObjectID(),
		Standar:          1, // Set standar sesuai kebutuhan (dapat dari request)
		Deskripsi:        "Deskripsi dokumen",
		LinkDokumen:      fileURL,
		Status:           "Penetapan",        // Bagian dari PPEPP
		CreatedBy:        "user@example.com", // Ambil dari JWT/Context
		LastModifiedBy:   "user@example.com",
		LastModifiedDate: time.Now().Unix(),
	}

	// Simpan akreditasi ke MongoDB
	// ...
}
