package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sistemakreditasi/backend-akreditasi/config"
	"github.com/sistemakreditasi/backend-akreditasi/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// UploadPDF handles the uploading of a PDF file to Google Drive and saves metadata to MongoDB
func UploadPDF(w http.ResponseWriter, r *http.Request) {
	// Parse the form to retrieve the uploaded file
	err := r.ParseMultipartForm(10 << 20) // limit file size to 10MB
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Get file from the form
	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Connect to Google Drive using environment variable for credentials
	ctx := context.Background()
	credentials := os.Getenv("GOOGLE_CREDENTIALS")
	if credentials == "" {
		http.Error(w, "Google credentials are missing", http.StatusInternalServerError)
		return
	}

	// Initialize Google Drive service
	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(credentials)))
	if err != nil {
		fmt.Printf("Google Drive connection error: %v\n", err)
		http.Error(w, "Unable to connect to Google Drive", http.StatusInternalServerError)
		return
	}

	// Create a new file in Google Drive
	driveFile := &drive.File{Name: handler.Filename, MimeType: "application/pdf"}
	uploadedFile, err := driveService.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		fmt.Printf("Error uploading to Google Drive: %v\n", err)
		http.Error(w, "Failed to upload file to Google Drive", http.StatusInternalServerError)
		return
	}

	// Save PDF metadata to MongoDB
	pdfDocument := models.PDFDocument{
		ID:         primitive.NewObjectID(),
		FileName:   handler.Filename,
		FileID:     uploadedFile.Id,
		UploadedBy: "user@example.com", // Replace with actual user, can be taken from JWT or request context
		UploadedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	// Access MongoDB collection
	collection := config.Mongoconn.Collection("pdf_documents")

	_, err = collection.InsertOne(ctx, pdfDocument)
	if err != nil {
		fmt.Printf("Error saving metadata to MongoDB: %v\n", err)
		http.Error(w, "Failed to save document metadata", http.StatusInternalServerError)
		return
	}

	// Respond with the uploaded file metadata
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pdfDocument)
}

// DownloadPDF handles the downloading of a PDF file from Google Drive
func DownloadPDF(w http.ResponseWriter, r *http.Request) {
	// Get the file ID from the URL parameters
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		http.Error(w, "Missing file_id parameter", http.StatusBadRequest)
		return
	}

	// Connect to Google Drive using environment variable for credentials
	ctx := context.Background()
	credentials := os.Getenv("GOOGLE_CREDENTIALS")
	if credentials == "" {
		http.Error(w, "Google credentials are missing", http.StatusInternalServerError)
		return
	}

	// Initialize Google Drive service
	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(credentials)))
	if err != nil {
		fmt.Printf("Google Drive connection error: %v\n", err)
		http.Error(w, "Unable to connect to Google Drive", http.StatusInternalServerError)
		return
	}

	// Retrieve the file from Google Drive
	resp, err := driveService.Files.Get(fileID).Download()
	if err != nil {
		fmt.Printf("Error downloading from Google Drive: %v\n", err)
		http.Error(w, "Failed to download file from Google Drive", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fileID))
	w.Header().Set("Content-Type", "application/pdf")

	// Stream the file content to the client
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		fmt.Printf("Error sending file to client: %v\n", err)
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
	}
}
