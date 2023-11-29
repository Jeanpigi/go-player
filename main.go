package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	musicFiles []string
	current    int
)

func loadMusicFiles(folder string) error {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			musicFiles = append(musicFiles, path)
		}
		return nil
	})
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	if len(musicFiles) == 0 {
		http.Error(w, "No music files available", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(musicFiles[current])
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Imprimir qué música está sonando
	fmt.Println("Playing:", musicFiles[current])

	w.Header().Set("Content-Type", "audio/mpeg")
	io.Copy(w, file)

	current = (current + 1) % len(musicFiles)
}

func main() {
	musicFolder := "./music"
	err := loadMusicFiles(musicFolder)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/stream", streamHandler)
	// Mostrar la dirección del servidor
	fmt.Println("Streaming server started on http://localhost:3006")
	http.ListenAndServe(":3006", nil)
}
