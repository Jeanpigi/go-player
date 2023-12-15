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
	musicFiles  []string
	playlist    []string
	currentSong int
)

func loadMusicFiles(folder string) error {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			musicFiles = append(musicFiles, path)
		}
		return nil
	})
}

func createPlaylist(files []string) {
	playlist = make([]string, len(files))
	copy(playlist, files)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(playlist), func(i, j int) {
		playlist[i], playlist[j] = playlist[j], playlist[i]
	})
}

func nextSong() string {
	if currentSong >= len(playlist) {
		createPlaylist(musicFiles) // Recreate the playlist once all songs have been played
		currentSong = 0
	}
	song := playlist[currentSong]
	currentSong++
	return song
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	if len(musicFiles) == 0 {
		http.Error(w, "No music files available", http.StatusInternalServerError)
		return
	}

	songPath := nextSong()
	file, err := os.Open(songPath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	w.Header().Set("Content-Disposition", "inline; filename=\""+path.Base(songPath)+"\"")
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "no-cache")

	fmt.Println("Playing:", songPath)

	io.Copy(w, file)
}

func main() {
	musicFolder := "./music"
	err := loadMusicFiles(musicFolder)
	if err != nil {
		panic(err)
	}

	createPlaylist(musicFiles)

	http.HandleFunc("/stream", streamHandler)

	fmt.Println("Streaming server started on http://localhost:3006")
	err = http.ListenAndServe(":3006", nil)
	if err != nil {
		panic(err)
	}
}
