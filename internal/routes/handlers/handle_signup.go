package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/internal/services/auth/rt"
	"github.com/Laelapa/PlateOps/util/net"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupResponse struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message,omitempty"`
	UserID       uuid.UUID `json:"user_id,omitempty"`
	JWT          string    `json:"jwt,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (h *Handler) HandlePostSignup(w http.ResponseWriter, r *http.Request) {

	var rBody signupRequest
	ctx := r.Context()

	// TODO: Consider deescalating log level.
	h.logger.LogRequestInfo("Processing signup request", r)

	// Decode the request body.
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.LogRequestWarn("Couldn't decode request body", r)
		h.logger.LogAppWarn("Couldn't decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the contents.
	if err := validateSignupRequest(rBody); err != nil {
		h.logger.LogAppInfo("Invalid signup request contents", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check username availability.
	// TODO: Consider spinning it of into a separate util function for code clarity.
	if _, err := h.queries.GetUserByUsername(ctx, rBody.Username); err == nil {
		h.logger.LogAppInfo(
			"Username already taken",
			zap.String("request_username", rBody.Username),
			zap.Error(err),
		)
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	} else if !errors.Is(err, pgx.ErrNoRows) { // Check for errors other than `pgx.ErrNoRows`.
		// If the error WAS `pgx.ErrNoRows`, it would mean the username is available.
		h.logger.LogAppError("Error checking username", err, zap.String("username", rBody.Username))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check email availability.
	// TODO: Consider spinning it of into a separate util function for code clarity.
	if _, err := h.queries.GetUserByEmail(ctx, rBody.Email); err == nil {
		h.logger.LogAppInfo(
			"Email already taken",
			zap.String("request_email", rBody.Email),
			zap.Error(err),
		)
		http.Error(w, "Email already taken", http.StatusConflict)
		return
	} else if !errors.Is(err, pgx.ErrNoRows) { // Similar to username flow.
		h.logger.LogAppError("Error checking email", err, zap.String("email", rBody.Email))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Hash the password.
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(rBody.Password), 14)
	if err != nil {
		h.logger.LogAppError("Error hashing password", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the user.
	createParams := repository.CreateUserParams{
		Username:     rBody.Username,
		Email:        rBody.Email,
		PasswordHash: string(hashedPwd),
	}

	user, err := h.queries.CreateUser(ctx, createParams)
	if err != nil {
		h.logger.LogAppError("Failed to create user", err,
			zap.String("username", rBody.Username),
			zap.String("email", rBody.Email),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate and register auth tokens.

	// Generate the refresh token.
	rToken, rtExpiresAt, err := h.tokenAuthority.IssueRT()
	if err != nil {
		h.logger.LogAppError("Failed to issue refresh token", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Register the refresh token.
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

	// Generate the JWT.
	jwt, err := h.tokenAuthority.IssueJWT(user.ID)
	if err != nil {
		h.logger.LogAppError("Failed to issue JWT", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the response.
	resp := signupResponse{
		Success:      true,
		Message:      "User created successfully",
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

	// TODO: Consider deescalating log level.
	h.logger.LogRequestInfo("Signup request processed successfully", r)

}

// TODO: Also check for zero values, indicating wrong json fields in request.
func validateSignupRequest(rBody signupRequest) error {

	// TODO: Validate username, length, regex

	// TODO: Validate email, length, regex

	// TODO: Validate password, length (max 72 for bcrypt)

	return nil
}
