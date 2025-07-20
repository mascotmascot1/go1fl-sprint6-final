package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// cleanupRequest closes and discards the request body.
// It is used to ensure proper cleanup of request resources.
func cleanupRequest(r *http.Request) {
	if r != nil {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}
}

// RootHandle handles the root endpoint "/".
// It reads and returns the contents of the "index.html" file as an HTML page.
//
// On error, it returns http.StatusInternalServerError.
func RootHandle(w http.ResponseWriter, r *http.Request) {
	defer cleanupRequest(r)

	data, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("error opening/reading file: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UploadConvertHandle handles the "/upload" endpoint.
//
// It performs the following steps:
// - Parses the multipart form and ensures exactly one uploaded file is provided.
// - Reads the contents of the uploaded file.
// - Passes the contents to service.Convert to get the converted representation.
// - Creates a new local file in the "output" directory with the result.
// - Returns the converted string as a plain text response.
//
// Returns http.StatusBadRequest if the user input is invalid (e.g., no file uploaded,
// more than one file uploaded, or conversion error), and http.StatusInternalServerError
// for all other internal errors.
func UploadConvertHandle(w http.ResponseWriter, r *http.Request) {
	const (
		formFileField = "myFile"
		outputDir     = "output"
	)

	defer cleanupRequest(r)

	if err := r.ParseMultipartForm(16 << 20); err != nil {
		http.Error(w, fmt.Sprintf("error multipart parsing: %s", err), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File[formFileField]
	if len(files) != 1 {
		http.Error(w, "exactly one file must be uploaded", http.StatusBadRequest)
		return
	}

	file, err := files[0].Open()
	if err != nil {
		http.Error(w, fmt.Sprintf("error opening file: %s", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading file: %s", err), http.StatusInternalServerError)
		return
	}

	convertedData, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("error converting file: %s", err), http.StatusBadRequest)
		return
	}

	if err = os.MkdirAll(outputDir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("error creating output directory: %s", err), http.StatusInternalServerError)
		return
	}

	outFile, err := os.Create(filepath.Join(outputDir, time.Now().UTC().Format("20060102_150405")+filepath.Ext(files[0].Filename)))
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating output file: %s", err), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err = outFile.Write([]byte(convertedData)); err != nil {
		http.Error(w, fmt.Sprintf("error writing file: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(convertedData))
}
