package controller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sistemakreditasi/backend-akreditasi/helper"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Fungsi untuk membuat service Google Drive
func getDriveService() (*drive.Service, error) {
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsFile("path/to/credentials.json"))
	if err != nil {
		return nil, fmt.Errorf("tidak dapat membuat service: %v", err) // Ubah huruf kapital menjadi huruf kecil
	}
	return srv, nil
}

// Fungsi untuk mengunggah file ke Google Drive
func UploadPDF(respw http.ResponseWriter, req *http.Request) {
	// Parse form untuk mendapatkan file yang diunggah
	err := req.ParseMultipartForm(10 << 20) // batas maksimal 10 MB
	if err != nil {
		http.Error(respw, "Gagal mem-parsing form", http.StatusBadRequest)
		return
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		http.Error(respw, "Tidak ada file yang diunggah", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Simpan file sementara ke server
	savePath := filepath.Join(os.TempDir(), handler.Filename)
	tempFile, err := os.Create(savePath)
	if err != nil {
		http.Error(respw, "Gagal menyimpan file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Salin konten file yang diunggah ke file sementara
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(respw, "Gagal menyimpan file", http.StatusInternalServerError)
		return
	}

	// Membuat service Google Drive
	srv, err := getDriveService()
	if err != nil {
		http.Error(respw, "Gagal terhubung ke Google Drive", http.StatusInternalServerError)
		return
	}

	// Membuka file yang diunggah untuk dibaca
	fileUpload, err := os.Open(savePath)
	if err != nil {
		http.Error(respw, "Gagal membuka file", http.StatusInternalServerError)
		return
	}
	defer fileUpload.Close()

	// Membuat metadata file untuk Drive
	driveFile := &drive.File{
		Name:    handler.Filename,
		Parents: []string{"16H3QVeFv3lQ-VjxEfCEnuG_fKqwYE11R"}, // Sesuaikan ID folder di Google Drive
	}

	// Unggah file ke Google Drive
	uploadedFile, err := srv.Files.Create(driveFile).Media(fileUpload).Do()
	if err != nil {
		http.Error(respw, "Gagal mengunggah file", http.StatusInternalServerError)
		return
	}

	// Hapus file sementara dari server
	os.Remove(savePath)

	// Response sukses
	resp := map[string]interface{}{
		"message": "File berhasil diunggah",
		"file_id": uploadedFile.Id,
	}
	helper.WriteJSON(respw, http.StatusOK, resp)
}

// Fungsi untuk mendownload file dari Google Drive
func DownloadPDF(respw http.ResponseWriter, req *http.Request) {
	fileId := req.URL.Query().Get("fileId") // Ambil fileId dari query parameter

	// Membuat service Google Drive
	srv, err := getDriveService()
	if err != nil {
		http.Error(respw, "Gagal terhubung ke Google Drive", http.StatusInternalServerError)
		return
	}

	// Mendapatkan file dari Google Drive
	res, err := srv.Files.Get(fileId).Download()
	if err != nil {
		http.Error(respw, "Gagal mengunduh file", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Mengirim file sebagai respons
	respw.Header().Set("Content-Disposition", "attachment; filename="+fileId)
	respw.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	respw.WriteHeader(http.StatusOK)
	io.Copy(respw, res.Body)
}
