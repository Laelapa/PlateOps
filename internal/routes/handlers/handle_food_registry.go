package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Laelapa/PlateOps/util/validate"
	"go.uber.org/zap"
)

type requestCreateFood struct {
	Name                   string  `json:"name"`
	Gtin                   string  `json:"gtin,omitempty"`
	Category               string  `json:"category,omitempty"`
	Description            string  `json:"description,omitempty"`
	UnitType               string  `json:"unit_type"` // Type of unit, e.g., "grams", "liters", "items", "portions" etc.
	Quantity               int     `json:"quantity"`
	PortionCount           int     `json:"portion_count"`                      // Number of portions in the food item
	ExpirationAfterOpening int     `json:"expiration_after_opening,omitempty"` // in days
	NutrientsPerItem       bool    `json:"nutrients_per_item,omitempty"`       // true: per item/portion, false: per 100g
	Calories               float32 `json:"calories,omitempty"`
	Fats                   float32 `json:"fats,omitempty"`
	Saturated              float32 `json:"saturated,omitempty"`
	Carbs                  float32 `json:"carbs,omitempty"`
	Sugars                 float32 `json:"sugars,omitempty"`
	Protein                float32 `json:"protein,omitempty"`
	Fiber                  float32 `json:"fiber,omitempty"`
	Sodium                 float32 `json:"sodium,omitempty"`
}

type responseCreateFood struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	ProductID int32  `json:"product_id,omitempty"`
}

func (h *Handler) HandlePostFood(w http.ResponseWriter, r *http.Request) {

	var rBody requestCreateFood
	ctx := r.Context()

	// Extract user ID from context - this should be set by the authentication middleware
	userID := r.Context().Value("userID")
	// TODO: Consider edge cases, check if they are accounted for in the middleware

	h.logger.LogRequestInfo("Inserting new food to the registry", r)

	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.LogRequestWarn("Couldn't decode request body", r)
		h.logger.LogAppWarn("Couldn't decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if rBody.Name == "" || rBody.UnitType == "" || rBody.Quantity <= 0 {
		h.logger.LogRequestWarn("Invalid request body: missing required fields", r)
		http.Error(w, "Bad request: missing required fields", http.StatusBadRequest)
		return
	}

	// Check if `Name` is already used
	existingFood, err := h.queries.GetFoodEntriesByName(ctx, rBody.Name)
	if err != nil {
		h.logger.LogAppError("Failed to check existing food entries", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(existingFood) > 0 {
		h.logger.LogRequestWarn("Food with the same name already exists", r)
		http.Error(w, "Food with the same name already exists", http.StatusConflict)
		return
	}

	// Convert request fields to database parameters
	dbParams, err := convertToCreateFoodEntryParams(rBody, userID)
	if err != nil {
		h.logger.LogAppError("Error converting request to database params", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Register the food entry with the database
	foodEntry, err := h.queries.CreateFoodEntry(ctx, dbParams)
	if err != nil {
		h.logger.LogAppError("Failed to create food entry in the database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create success response
	resp := responseCreateFood{
		Success:   true,
		Message:   "Food entry created successfully",
		ProductID: foodEntry.ProductID,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.LogAppError("Failed to encode response", err)
	}

	h.logger.LogAppInfo(
		"Food entry created successfully",
		zap.Int32("product_id", foodEntry.ProductID),
		zap.String("user_id", userID.(string)),
	)
}

func validateCreateFoodRequest(req requestCreateFood) error {
	if err := validate.StringRequired(req.Name); err != nil {
		return err
	}

	if req.Gtin != "" { // if not omitted
		if err := validate.GTIN(req.Gtin); err != nil {
			return err
		}
	}

	if err := validate.UnitType(req.UnitType); err != nil {
		return err
	}

	if req.Category != "" {
		if err := validate.String(req.Category); err != nil {
			return err
		}
	}

	if req.Description != "" {
		if err := validate.Text(req.Description); err != nil {
			return err
		}
	}

	if req.Quantity != 0 {
		if err := validate.Positive()
	}
}
