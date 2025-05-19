package handlers

import (
	"github.com/Laelapa/GoHome/logging"
	"github.com/Laelapa/PlateOps/auth/tokenauthority"
	"github.com/Laelapa/PlateOps/internal/repository"
)

type Handler struct {
	logger         *logging.Logger
	queries        *repository.Queries
	tokenAuthority *tokenauthority.TokenAuthority
}

func New(logger *logging.Logger, queries *repository.Queries, tokenAuthority *tokenauthority.TokenAuthority) *Handler {
	return &Handler{
		logger:         logger,
		queries:        queries,
		tokenAuthority: tokenAuthority,
	}
}
