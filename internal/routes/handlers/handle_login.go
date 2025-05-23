package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/PlateOps/internal/services/auth/rt"
	"github.com/Laelapa/PlateOps/util/net"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message,omitempty"`
	UserID       uuid.UUID `json:"user_id,omitempty"`
	JWT          string    `json:"jwt,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (h *Handler) HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	var rBody loginRequest
	ctx := r.Context()

	h.logger.LogRequestInfo("Processing login request", r)

	// Decode the request body.
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.LogRequestWarn("Couldn't decode request body", r)
		h.logger.LogAppWarn("Couldn't decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the contents.
	if err := validateLoginRequest(rBody); err != nil {
		h.logger.LogAppInfo("Invalid login request contents", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Fetch the user.
	user, err := h.queries.GetUserByUsername(ctx, rBody.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.logger.LogAppInfo("User not found", zap.String("request_username", rBody.Username))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.logger.LogAppError("Couldn't fetch user from database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Verify the password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(rBody.Password)); err != nil {
		h.logger.LogAppInfo("Invalid password", zap.String("request_username", rBody.Username))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate and register auth tokens

	// Issue the JWT
	jwt, err := h.tokenAuthority.IssueJWT(user.ID)
	if err != nil {
		h.logger.LogAppError("Couldn't generate JWT", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Issue the RT
	rToken, rtExpiresAt, err := h.tokenAuthority.IssueRT()
	if err != nil {
		h.logger.LogAppError("Couldn't generate refresh token", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Register the refresh token
	rtParams := rt.Params{
		UserID:    user.ID,
		Token:     rToken,
		ExpiresAt: rtExpiresAt,
		UserAgent: r.UserAgent(),
		IPAddress: net.GetFlyClientIP(r),
	}
	if err := rt.RegisterNewToken(h.queries, h.logger, rtParams); err != nil {
		h.logger.LogAppError("Failed to register refresh token", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build the response
	resp := loginResponse{
		Success:      true,
		Message:      "Login successful",
		UserID:       user.ID,
		JWT:          jwt,
		RefreshToken: rToken,
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

	h.logger.LogRequestInfo("Login request processed successfully", r)
}

// TODO: Also check for zero values, indicating wrong json fields in request.
func validateLoginRequest(rBody loginRequest) error {

	// TODO: Validate username, length, regex

	// TODO: Validate password, length (max 72 for bcrypt)

	return nil
}
