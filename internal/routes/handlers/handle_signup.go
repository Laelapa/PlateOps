package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/PlateOps/internal/repository"

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
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	JWT          string `json:"jwt,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (h *Handler) HandlePostSignup(w http.ResponseWriter, r *http.Request) {

	var rBody signupRequest
	ctx := r.Context()

	h.logger.LogRequestInfo("Processing signup request", r)

	// Decode the request body.
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.LogRequestError("Error decoding request body", r, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the contents.
	if err := validateSignupRequest(rBody); err != nil {
		h.logger.LogRequestError("Invalid signup request", r, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check username availability.
	if _, err := h.queries.GetUserByUsername(ctx, rBody.Username); err == nil {
		h.logger.LogRequestError("Username already taken", r, err)
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	} else if !errors.Is(err, pgx.ErrNoRows) { // Check for errors other than `pgx.ErrNoRows`.
		// If the error WAS `pgx.ErrNoRows`, it would mean the username is available.
		h.logger.LogAppError("Error checking username", err, zap.String("username", rBody.Username))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check email availability.
	if _, err := h.queries.GetUserByEmail(ctx, rBody.Email); err == nil {
		h.logger.LogRequestError("Email already taken", r, err)
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

	// Respond with the success message and user tokens.

}

func validateSignupRequest(rBody signupRequest) error {

	// TODO: Validate username, length, regex

	// TODO: Validate email, length, regex

	// TODO: Validate password, length (max 72 for bcrypt)

	return nil
}
