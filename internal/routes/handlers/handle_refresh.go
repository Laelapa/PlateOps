package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Laelapa/PlateOps/util/auth"
	"github.com/Laelapa/PlateOps/util/net"
	"github.com/Laelapa/PlateOps/util/validate"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
		h.logger.LogRequestWarn("Couldn't decode request body", r)
		h.logger.LogAppWarn("Couldn't decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the contents.
	if err := validateRefreshRequest(
		rBody, 
		h.tokenAuthority.GetRefreshTokenSizeInBytes(),
	); err != nil {
		h.logger.LogAppInfo("Invalid token refresh request", zap.Error(err), zap.String("request_user_id", rBody.UserID.String()))
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

		h.logger.LogAppWarn(
			"Token verification failed",
			zap.String("request_refresh_token", rBody.RefreshToken),
			zap.String("request_user_id", rBody.UserID.String()),
			zap.String("request_user_agent", userAgent),
			zap.String("request_ip_address", ipAddress),
			zap.Error(err))
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
	if _, err := w.Write(respMarshalled); err != nil {
		h.logger.LogAppError("Failed to write response", err)
	}

	h.logger.LogRequestInfo("JWT refreshed successfully", r)

}

func validateRefreshRequest(rBody refreshRequest, rtSizeInBytes int) error {

	if err := validate.UUID(rBody.UserID); err != nil {
		return err
	}

	if err := validate.RefreshToken(rBody.RefreshToken, rtSizeInBytes); err != nil {
		return err
	}

	return nil
}
