package handlers

import (
	"net/http"

	"github.com/Laelapa/GoHome/internal/interface/templates"
)

func (h *Handler) HandleGetStack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// TODO: pull the parameters from env/secrets & request to pass them to the template instead of hardcoding them
	if err := templates.Stack("Your website's name - Tech Stack", "example.com", "/stack").Render(r.Context(), w); err != nil {
		h.Logger.LogRequestError("Failed to render stack page", r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.Logger.LogRequestInfo("Rendered: Stack page", r)
}
