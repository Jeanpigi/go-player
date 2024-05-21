package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/jeanpigi/go-player/internal/playlist"
)

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	if len(playlist.Playlist) == 0 {
		http.Error(w, "No music files available", http.StatusInternalServerError)
		return
	}

	songPath := playlist.NextSong()
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
