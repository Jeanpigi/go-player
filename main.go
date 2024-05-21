package main

import (
	"fmt"
	"net/http"

	"github.com/jeanpigi/go-player/internal/handlers"
	"github.com/jeanpigi/go-player/internal/music"
	"github.com/jeanpigi/go-player/internal/playlist"
)

func main() {
	musicFolder := "./music"
	err := music.LoadMusicFiles(musicFolder)
	if err != nil {
		panic(err)
	}

	playlist.CreatePlaylist()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/stream", handlers.StreamHandler)
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	fmt.Println("Streaming server started on http://localhost:3006")
	err = http.ListenAndServe(":3006", nil)
	if err != nil {
		panic(err)
	}
}
