package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Laelapa/PlateOps/internal/repository"
	"go.uber.org/zap"
)

func (h *Handler) HandlePostFood(w http.ResponseWriter, r *http.Request) {

	var rBody repository.CreateFoodEntryParams
	ctx := r.Context()

	h.logger.LogRequestInfo("Inserting new food to the registry", r)

	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.LogRequestWarn("Couldn't decode request body", r)
		h.logger.LogAppWarn("Couldn't decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err := h.queries.CreateFoodEntry(ctx, rBody)
	if err != nil {
		h.logger.LogAppError("Couldn't insert food entry", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
