package handlers

import (
	"net/http"

	"github.com/jeanpigi/go-player/internal/templates"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templates.RenderTemplate(w, "web/templates/home.html", nil)
}
