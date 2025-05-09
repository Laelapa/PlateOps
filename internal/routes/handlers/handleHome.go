package handlers

import (
	"net/http"

	"github.com/Laelapa/GoHome/internal/interface/templates"
)

func (h *Handler) HandleGetHome(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	// TODO: pull the parameters from env/secrets & request to pass them to the template after refactoring the templates
	if err := templates.Home().Render(r.Context(), w); err != nil {
		h.Logger.LogRequestError("Failed to render home page", r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return

	}

	h.Logger.LogRequestInfo("Rendered: Home page", r)
}
