package handlers

import (
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
