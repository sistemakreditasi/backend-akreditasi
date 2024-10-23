package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sistemakreditasi/backend-akreditasi/helper"
	"github.com/sistemakreditasi/backend-akreditasi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func UploadPDF(w http.ResponseWriter, r *http.Request) {
	// Parse form to retrieve uploaded file
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "File terlalu besar", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Google Drive credentials
	credentials := os.Getenv("GOOGLE_CREDENTIALS")
	if credentials == "" {
		http.Error(w, "Google credentials are missing", http.StatusInternalServerError)
		return
	}

	// Connect to Google Drive
	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(credentials)))
	if err != nil {
		http.Error(w, "Unable to connect to Google Drive", http.StatusInternalServerError)
		return
	}

	// Upload file to Google Drive
	driveFile := &drive.File{Name: handler.Filename, MimeType: "application/pdf"}
	uploadedFile, err := driveService.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		http.Error(w, "Failed to upload file to Google Drive", http.StatusInternalServerError)
		return
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGOSTRING")))
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database("db_akreditasi").Collection("pdf_documents")

	// Save PDF metadata to MongoDB
	pdfDocument := model.PDFDocument{
		ID:         primitive.NewObjectID(),
		FileName:   handler.Filename,
		FileID:     uploadedFile.Id,
		UploadedBy: "user@example.com", // Bisa diganti dengan user dari JWT/token
		UploadedAt: time.Now(),
	}

	_, err = collection.InsertOne(ctx, pdfDocument)
	if err != nil {
		http.Error(w, "Failed to save document metadata", http.StatusInternalServerError)
		return
	}

	// Respond with uploaded file metadata
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pdfDocument)
}

// DownloadPDF handles the downloading of a PDF file from Google Drive
func DownloadPDF(w http.ResponseWriter, r *http.Request) {
	// Get the file ID from the URL parameters
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		helper.ErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "Missing file_id parameter")
		return
	}

	// Connect to Google Drive using environment variable for credentials
	ctx := context.Background()
	credentials := os.Getenv("GOOGLE_CREDENTIALS")
	if credentials == "" {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "Google credentials are missing")
		return
	}

	// Initialize Google Drive service
	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(credentials)))
	if err != nil {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "Unable to connect to Google Drive")
		return
	}

	// Retrieve the file from Google Drive
	resp, err := driveService.Files.Get(fileID).Download()
	if err != nil {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "Failed to download file from Google Drive")
		return
	}
	defer resp.Body.Close()

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fileID))
	w.Header().Set("Content-Type", "application/pdf")

	// Stream the file content to the client
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "Failed to send file")
	}
}
