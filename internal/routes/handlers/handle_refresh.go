package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Laelapa/PlateOps/util/auth"
	"github.com/Laelapa/PlateOps/util/net"
	"github.com/google/uuid"
)

type refreshRequest struct {
	RefreshToken string    `json:"refresh_token"`
	UserID       uuid.UUID `json:"user_id"`
	// TODO: [Very low priority] Consider implementing a timestamp field for security (identify delayed requests / clock drift etc.)
}

type refreshResponse struct {
	JWT string `json:"jwt,omitempty"`
}

func (h *Handler) HandlePostRefresh(w http.ResponseWriter, r *http.Request) {

	var rBody refreshRequest
	ctx := r.Context()

	// Fetch relevant security headers.
	userAgent := r.UserAgent()
	ipAddress := net.GetFlyClientIP(r)

	h.logger.LogRequestInfo("Processing token refresh request", r)

	// Decode the request body.
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.LogRequestError("Error decoding request body", r, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the contents.
	if err := validateRefreshRequest(rBody); err != nil {
		h.logger.LogRequestError("Invalid token refresh request", r, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Verify the request info & fingerprints against the database
	if err := auth.VerifyRefreshToken(
		ctx,
		h.queries,
		rBody.RefreshToken,
		rBody.UserID,
		userAgent,
		ipAddress,
		time.Now().UTC(),
	); err != nil {
		if errors.Is(err, auth.ErrDatabaseFailure) {
			h.logger.LogAppError("Database error during token verification", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.logger.LogRequestError("Token verification failed", r, err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate a fresh JWT
	jwt, err := h.tokenAuthority.IssueJWT(rBody.UserID)
	if err != nil {
		h.logger.LogAppError("Failed to issue JWT", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the response.
	resp := refreshResponse{
		JWT: jwt,
	}

	respMarshalled, err := json.Marshal(resp)
	if err != nil {
		h.logger.LogAppError("Failed to marshal response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respMarshalled)

	h.logger.LogRequestInfo("JWT refreshed successfully", r)

}

// validateRefreshRequest verifies that the contents of the received request are valid as far as allowed symbols and length are concerned.
func validateRefreshRequest(rBody refreshRequest) error {

	// TODO: Validate username, length, regex

	// TODO: Validate refresh token, length, regex

	return nil
}
