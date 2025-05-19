package handlers

import "net/http"

func (h *Handler) HandleGetHealth(w http.ResponseWriter, r *http.Request) {

	h.logger.LogRequestInfo("Rendered: Health check", r)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		h.logger.LogAppError("Error writing response", err)
	}
}
