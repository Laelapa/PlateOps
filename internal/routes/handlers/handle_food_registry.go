package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Laelapa/PlateOps/internal/repository"
	"go.uber.org/zap"
)

type createFoodRequest struct {
	Name                   string  `json:"name"`
	Gtin                   string  `json:"gtin,omitempty"`
	Category               string  `json:"category,omitempty"`
	Description            string  `json:"description,omitempty"`
	UnitType               string  `json:"unit_type"`	// Type of unit, e.g., "grams", "liters", "items", "portions" etc.
	Quantity               int     `json:"quantity"`
	PortionCount           int     `json:"portion_count"`	// Number of portions in the food item
	ExpirationAfterOpening int     `json:"expiration_after_opening,omitempty"` // in days
	NutrientsPerItem       bool    `json:"nutrients_per_item,omitempty"`	// true: per item/portion, false: per 100g
	Calories               float64 `json:"calories,omitempty"`
	Fats                   float64 `json:"fats,omitempty"`	
	Saturated              float64 `json:"saturated,omitempty"` 
	Carbs                  float64 `json:"carbs,omitempty"`
	Sugars                 float64 `json:"sugars,omitempty"`
	Protein                float64 `json:"protein,omitempty"`
	Fiber                  float64 `json:"fiber,omitempty"`
	Sodium                 float64 `json:"sodium,omitempty"`
}

type createFoodResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	ProductID string `json:"product_id,omitempty"`

func (h *Handler) HandlePostFood(w http.ResponseWriter, r *http.Request) {

	var rBody createFoodRequest
	ctx := r.Context()

	// Extract user ID from context - this should be set by the authentication middleware
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok ||  {

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
