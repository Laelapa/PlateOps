package handlers

import (
	"github.com/Laelapa/GoHome/logging"

	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/guarddoggo"
)

// TODO: Don't export the fields, set up a constructor function
type Handler struct {
	Logger  *logging.Logger
	Queries *repository.Queries
	Azor    guarddoggo.GuardDoggo
}
