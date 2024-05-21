package handlers

import (
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jeanpigi/go-player/internal/music"
)

func isMusicFile(header *multipart.FileHeader) bool {
	mimeType := header.Header.Get("Content-Type")
	return mimeType == "audio/mpeg"
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/upload.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		data := struct {
			Title string
		}{
			Title: "Upload Music",
		}
		tmpl.ExecuteTemplate(w, "layout", data)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20) // maxMemory 10MB
		file, handler, err := r.FormFile("musicFile")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if !isMusicFile(handler) {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		dst, err := os.Create(filepath.Join("./music", handler.Filename))
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		music.MusicFiles = append(music.MusicFiles, dst.Name())

		fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
	}
}
