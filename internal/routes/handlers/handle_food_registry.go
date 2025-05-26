package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/util/typeconvert"
	"github.com/Laelapa/PlateOps/util/validate"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type requestCreateFood struct {
	Name                   string  `json:"name"`
	Gtin                   string  `json:"gtin,omitempty"`
	Category               string  `json:"category,omitempty"`
	Description            string  `json:"description,omitempty"`
	UnitType               string  `json:"unit_type"` // Type of unit, e.g., "grams", "liters", "items", "portions" etc.
	Quantity               int     `json:"quantity"`
	PortionCount           int     `json:"portion_count,omitempty"`                      // Number of portions in the food item
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
	// TODO: Consider edge cases, check if they are accounted for in the middleware
	userID, ok := ctx.Value("userID").(pgtype.UUID)
	if !ok {
		h.logger.LogRequestWarn("Invalid or missing user ID in context", r)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	// Validate request fields
	if err := validateCreateFoodRequest(rBody); err != nil {
		h.logger.LogRequestWarn("Invalid request body", r)
		h.logger.LogAppWarn("Invalid request body", zap.Error(err))
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
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
		zap.String("user_id", typeconvert.PgtypeUUIDToString(userID)),
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
		if err := validate.Positive(req.Quantity); err != nil {
			return err
		}
	}

	if req.PortionCount == 0 {
		req.PortionCount = 1
	} else {
		if err := validate.Positive(req.PortionCount); err != nil {
			return err
		}
	}

	if req.ExpirationAfterOpening != 0 {
		if err := validate.NonNegative(req.ExpirationAfterOpening); err != nil {
			return err
		}
	}

	for _, nutrient := range []float32{
		req.Calories,
		req.Fats,
		req.Saturated,
		req.Carbs,
		req.Sugars,
		req.Protein,
		req.Fiber,
		req.Sodium,
	} {
		if err := validate.NonNegative(nutrient); err != nil {
			return err
		}
	}

	return nil
}

func convertToCreateFoodEntryParams(req requestCreateFood, userID pgtype.UUID) (repository.CreateFoodEntryParams, error) {

	params := repository.CreateFoodEntryParams{
		Name:                   req.Name,
		Gtin:                   typeconvert.StringToPgtypeText(req.Gtin),
		Category:               typeconvert.StringToPgtypeText(req.Category),
		Description:            typeconvert.StringToPgtypeText(req.Description),
		UnitType:               req.UnitType,
		Quantity:               typeconvert.IntToPgtypeInt4(req.Quantity),
		PortionCount:           typeconvert.IntToPgtypeInt4(req.PortionCount),
		ExpirationAfterOpening: typeconvert.IntToPgtypeInt4(req.ExpirationAfterOpening),
		NutrientsPerItem:       typeconvert.BoolToPgtypeBool(req.NutrientsPerItem),
		Calories:               typeconvert.Float32ToPgtypeFloat4(req.Calories),
		Fats:                   req.Fats,
		Saturated:              req.Saturated,
		Carbs:                  req.Carbs,
		Sugars:                 req.Sugars,
		Protein:                req.Protein,
		Fiber:                  req.Fiber,
		Sodium:                 req.Sodium,
		CreatedBy:              userID,
	}


	return params, nil
}