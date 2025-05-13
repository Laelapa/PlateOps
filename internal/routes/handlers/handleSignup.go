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
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	// TODO: return token
}

func (h *Handler) HandlePostSignup(w http.ResponseWriter, r *http.Request) {

	var rBody signupRequest
	ctx := r.Context()

	h.Logger.LogRequestInfo("Processing signup request", r)

	// Decode the request body.
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.Logger.LogRequestError("Error decoding request body", r, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the contents.
	if err := validateSignupRequest(rBody); err != nil {
		h.Logger.LogRequestError("Invalid signup request", r, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check username availability.
	if _, err := h.Queries.GetUserByUsername(ctx, rBody.Username); err == nil {
		h.Logger.LogRequestError("Username already taken", r, err)
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		h.Logger.LogAppError("Error checking username", err, zap.String("username", rBody.Username))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check email availability.
	if _, err := h.Queries.GetUserByEmail(ctx, rBody.Email); err == nil {
		h.Logger.LogRequestError("Email already taken", r, err)
		http.Error(w, "Email already taken", http.StatusConflict)
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		h.Logger.LogAppError("Error checking email", err, zap.String("email", rBody.Email))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Hash the password.
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(rBody.Password), 14)
	if err != nil {
		h.Logger.LogAppError("Error hashing password", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the user.
	createParams := repository.CreateUserParams{
		Username:     rBody.Username,
		Email:        rBody.Email,
		PasswordHash: string(hashedPwd),
	}

	user, err := h.Queries.CreateUser(ctx, createParams)
	if err != nil {
		h.Logger.LogAppError("Failed to create user", err,
			zap.String("username", rBody.Username),
			zap.String("email", rBody.Email),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with the created user.

}

func validateSignupRequest(rBody signupRequest) error {

	// TODO: Validate username, length, regex

	// TODO: Validate email, length, regex

	// TODO: Validate password, length (max 72 for bcrypt)

	return nil
}
