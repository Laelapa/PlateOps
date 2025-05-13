package handlers

import (
	"github.com/Laelapa/GoHome/logging"
	"github.com/Laelapa/PlateOps/internal/repository"
)

type Handler struct {
	Logger  *logging.Logger
	Queries *repository.Queries
}
