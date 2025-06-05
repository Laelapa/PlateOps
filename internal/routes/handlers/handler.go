package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Laelapa/PlateOps/auth/tokenauthority"
	"github.com/Laelapa/PlateOps/internal/repository"

	"github.com/Laelapa/GoHome/logging"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	logger         *logging.Logger
	queries        *repository.Queries
	tokenAuthority *tokenauthority.TokenAuthority
	validator      *validator.Validate
}

func New(logger *logging.Logger, queries *repository.Queries, tokenAuthority *tokenauthority.TokenAuthority) *Handler {
	return &Handler{
		logger:         logger,
		queries:        queries,
		tokenAuthority: tokenAuthority,
		validator:      validator.New(validator.WithRequiredStructEnabled()),
	}
}

// TODO: Future refactor go through multiple error types to slim down the number of parameters (by removing msg and status code)
func (h *Handler) HandleError(w http.ResponseWriter, r *http.Request, err error, msg string, statusCode int) {
	h.logger.LogRequestWarn(msg, r)
	h.logger.LogAppWarn(msg, zap.Error(err))
	http.Error(w, msg, statusCode)
}
