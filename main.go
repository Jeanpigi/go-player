package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

var (
	musicFiles []string
)

func loadMusicFiles(folder string) error {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Handle error in accessing the folder or file
		}
		if !info.IsDir() {
			musicFiles = append(musicFiles, path)
		}
		return nil
	})
}

func shuffleFiles(files []string) {
	rand.Seed(time.Now().UnixNano())
	for i := len(files) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		files[i], files[j] = files[j], files[i]
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	if len(musicFiles) == 0 {
		http.Error(w, "No music files available", http.StatusInternalServerError)
		return
	}

	// Select a random file
	randomIndex := rand.Intn(len(musicFiles))
	file, err := os.Open(musicFiles[randomIndex])
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Obtain file information
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	// Set HTTP headers
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	w.Header().Set("Content-Disposition", "inline; filename=\""+path.Base(musicFiles[randomIndex])+"\"")
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "no-cache")

	// Print what music is playing
	fmt.Println("Playing:", musicFiles[randomIndex])

	io.Copy(w, file)
}

func main() {
	musicFolder := "./music"
	err := loadMusicFiles(musicFolder)
	if err != nil {
		panic(err) // Handle error loading music files
	}

	shuffleFiles(musicFiles) // Shuffle music files

	http.HandleFunc("/stream", streamHandler)

	// Show server address
	fmt.Println("Streaming server started on http://localhost:3006")
	err = http.ListenAndServe(":3006", nil)
	if err != nil {
		panic(err) // Handle error starting the server
	}
}
