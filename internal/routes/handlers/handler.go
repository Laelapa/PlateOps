package handlers

import (
	"github.com/Laelapa/GoHome/logging"
	"github.com/Laelapa/PlateOps/internal/repository"
)

type Handler struct {
	logger  *logging.Logger
	queries *repository.Queries
}

func New(logger *logging.Logger, queries *repository.Queries) *Handler {
	return &Handler{
		logger:  logger,
		queries: queries,
	}
}
