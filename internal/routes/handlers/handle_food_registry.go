package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/util"
	"github.com/Laelapa/PlateOps/util/ctxutils"
	"github.com/Laelapa/PlateOps/util/parse"
	"github.com/Laelapa/PlateOps/util/typeconvert"
	"github.com/Laelapa/PlateOps/util/validate"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

var ErrGtinAlreadyExists = errors.New("GTIN already exists for another food entry")
var ErrQueryFailed = errors.New("failed to execute database query")

// requestFood represents the structure of the request body for creating or updating a food entry.
// All fields are set to be optional to accommodate the PATCH method.
// For fields that are necessary for other methods verify the request body in the handler.
type requestFood struct {
	Name                   string  `json:"name,omitempty"`
	Gtin                   string  `json:"gtin,omitempty"`
	Category               string  `json:"category,omitempty"`
	Description            string  `json:"description,omitempty"`
	UnitType               string  `json:"unit_type,omitempty"` // Type of unit, e.g., "grams", "liters", "items", "portions" etc.
	Quantity               int32   `json:"quantity,omitempty"`
	PortionCount           int32   `json:"portion_count,omitempty"`            // Number of portions in the food item
	ExpirationAfterOpening int32   `json:"expiration_after_opening,omitempty"` // in days
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

type responseFood struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	ProductID int32  `json:"product_id,omitempty"`
}

func (h *Handler) HandlePostFood(w http.ResponseWriter, r *http.Request) {

	var rBody requestFood
	ctx := r.Context()

	// Extract user ID from context - this should be set by the authentication middleware
	// TODO: Consider edge cases, check if they are accounted for in the middleware
	userID, ok := ctxutils.GetUserIDFromContext(ctx)
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

	// Verify that the required fields are present.
	if rBody.Name == "" || rBody.UnitType == "" || rBody.Quantity == 0 {
		h.logger.LogRequestWarn("Invalid request body: missing required fields", r)
		http.Error(w, "Bad request: missing required fields", http.StatusBadRequest)
		return
	}

	// Validate request fields
	if valErr := h.validateFoodRequestParams(ctx, rBody); valErr != nil {
		if errors.Is(valErr, ErrQueryFailed) {
			h.logger.LogAppError("Failed to validate request parameters", valErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		h.logger.LogRequestWarn("Invalid request body", r)
		h.logger.LogAppWarn("Invalid request body", zap.Error(valErr))
		http.Error(w, "Bad request: "+valErr.Error(), http.StatusBadRequest)
		return
	}

	// Convert request fields to database parameters
	dbParams := convertToCreateFoodEntryParams(rBody, userID)

	// Register the food entry with the database
	foodEntry, err := h.queries.CreateFoodEntry(ctx, dbParams)
	if err != nil {
		h.logger.LogAppError("Failed to create food entry in the database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create success response
	resp := responseFood{
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

func (h *Handler) HandlePatchFood(w http.ResponseWriter, r *http.Request) {

	var rBody requestFood
	ctx := r.Context()

	// TODO: Consider edge cases, check if they are accounted for in the middleware
	userID, ok := ctxutils.GetUserIDFromContext(ctx)
	if !ok {
		h.logger.LogRequestWarn("Invalid or missing user ID in context", r)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.logger.LogRequestInfo("Updating food entry in the registry", r)

	productID, err := parse.ID(r.PathValue("id"))
	if err != nil {
		h.logger.LogRequestWarn("Invalid product ID", r)
		h.logger.LogAppWarn("Invalid product ID", zap.Error(err))
		http.Error(w, "Bad request: invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.checkGtinUniqueness(ctx, rBody.Gtin, productID); err != nil {
		if errors.Is(err, ErrGtinAlreadyExists) {
			h.logger.LogRequestWarn("GTIN already exists for another food entry", r)
			http.Error(w, "GTIN already exists", http.StatusConflict)
			return
		}
		h.logger.LogAppError("Failed to check GTIN uniqueness", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.LogRequestWarn("Couldn't decode request body", r)
		h.logger.LogAppWarn("Couldn't decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// TODO: Consider locking editing to owner of the food entry
	// `if existingFood.CreatedBy != userID`

	// Validate request fields
	if valErr := h.validateFoodRequestParams(ctx, rBody); valErr != nil {
		if errors.Is(valErr, ErrQueryFailed) {
			h.logger.LogAppError("Failed to validate request parameters", valErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		h.logger.LogRequestWarn("Invalid request body", r)
		h.logger.LogAppWarn("Invalid request body", zap.Error(valErr))
		http.Error(w, "Bad request: "+valErr.Error(), http.StatusBadRequest)
		return
	}

}

// If a field is present, it checks its validity according to the business rules.
// Returns an error if any validation fails, or nil if all validations pass.
// This function does NOT verify the presence of required fields.
// It is expected that the handler will check for required fields before calling this function.
// 
// Parameters:
//   - ctx: The context for the request, used for database operations
//   - req: The request body containing the food entry data
//
// Returns:
//   - error: Returns nil if all validations pass, or an error if any validation fails.
//
// Possible errors:
//   - ErrQueryFailed: if there is an internal failure while executing the database query
//   - Propagated validation errors: if any of the fields fail validation checks
func (h *Handler) validateFoodRequestParams(ctx context.Context, req requestFood) error {

	if req.Gtin != "" { // if not omitted
		if err := validate.GTIN(req.Gtin); err != nil {
			return err
		}
		if err := h.checkGtinUniqueness(ctx, req.Gtin, -1); err != nil {
			return err
		}
	}

	if req.UnitType != "" {
		if err := validate.UnitType(req.UnitType); err != nil {
			return err
		}
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

	if req.PortionCount != 0 {
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

func convertToCreateFoodEntryParams(
	req requestFood,
	userID pgtype.UUID,
) repository.CreateFoodEntryParams {

	pcount := req.PortionCount
	if pcount == 0 {
		pcount = 1 // Default to 1 if not specified
	}

	params := repository.CreateFoodEntryParams{
		Name:                   req.Name,
		Gtin:                   typeconvert.StringToPgtypeText(req.Gtin),
		Category:               typeconvert.StringToPgtypeText(req.Category),
		Description:            typeconvert.StringToPgtypeText(req.Description),
		UnitType:               req.UnitType,
		Quantity:               typeconvert.Int32ToPgtypeInt4(req.Quantity),
		PortionCount:           typeconvert.Int32ToPgtypeInt4(pcount),
		ExpirationAfterOpening: typeconvert.Int32ToPgtypeInt4(req.ExpirationAfterOpening),
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

	return params
}

func (h *Handler) checkGtinUniqueness(ctx context.Context, gtin string, productID int32) error {

	// Check if GTIN is already used by another food entry
	existingFood, err := h.queries.GetFoodEntryByGtin(
		ctx,
		typeconvert.StringToPgtypeText(gtin),
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // GTIN is unique
		}
		return (&util.DeepError{BusinessErr: ErrQueryFailed, TechnicalErr: err})
	}

	if existingFood.ProductID != productID {
		return ErrGtinAlreadyExists
	}

	return nil // GTIN matches the current food entry
}
