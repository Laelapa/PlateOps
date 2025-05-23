package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func (h *Handler) HandlePostLogout(w http.ResponseWriter, r *http.Request) {
	var rBody logoutRequest
	ctx := r.Context()

	h.logger.LogRequestInfo("Processing logout request", r)

	// Decode the request body.
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.LogRequestWarn("Couldn't decode request body", r)
		h.logger.LogAppWarn("Couldn't decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the contents.
	if err := validateLogoutRequest(rBody); err != nil {
		h.logger.LogAppInfo("Invalid logout request contents", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Call the logout function.
	if err := h.queries.LogOutToken(ctx, rBody.RefreshToken); err != nil {
		h.logger.LogAppError("Logout failed", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send a success response.
	resp := logoutResponse{
		Success: true,
		Message: "Logout successful",
	}

	respMarshalled, err := json.Marshal(resp)
	if err != nil {
		h.logger.LogAppError("Couldn't marshal response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respMarshalled)

	h.logger.LogRequestInfo("Logout request processed successfully", r)

}

func validateLogoutRequest(rBody logoutRequest) error {

	return nil
}
