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
// Check for presence of a field in a request by checking if the pointer is nil.
type requestFood struct {
	Name                   *string  `json:"name,omitempty"                     validate:"omitempty,max=255"`
	Gtin                   *string  `json:"gtin,omitempty"                     validate:"omitempty,numeric, max=14"` // GTIN-14
	Category               *string  `json:"category,omitempty"                 validate:"omitempty,max=255"`
	Description            *string  `json:"description,omitempty"              validate:"omitempty,max=5000"`
	UnitType               *string  `json:"unit_type,omitempty"                validate:"omitempty,oneof=grams ml items portions"`
	Quantity               *int32   `json:"quantity,omitempty"`
	PortionCount           *int32   `json:"portion_count,omitempty"            validate:"omitempty,min=1"` // Number of portions in the food item
	ExpirationAfterOpening *int32   `json:"expiration_after_opening,omitempty" validate:"omitempty,min=0"` // in days
	NutrientsPerItem       *bool    `json:"nutrients_per_item,omitempty"`                                  // true: per item/portion, false: per 100g
	Calories               *float32 `json:"calories,omitempty"                 validate:"omitempty,min=0"`
	Fats                   *float32 `json:"fats,omitempty"                     validate:"omitempty,min=0"`
	Saturated              *float32 `json:"saturated,omitempty"                validate:"omitempty,min=0"`
	Carbs                  *float32 `json:"carbs,omitempty"                    validate:"omitempty,min=0"`
	Sugars                 *float32 `json:"sugars,omitempty"                   validate:"omitempty,min=0"`
	Protein                *float32 `json:"protein,omitempty"                  validate:"omitempty,min=0"`
	Fiber                  *float32 `json:"fiber,omitempty"                    validate:"omitempty,min=0"`
	Sodium                 *float32 `json:"sodium,omitempty"                   validate:"omitempty,min=0"`
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
		h.HandleError(w, r, err, "Couldn't decode request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(rBody); err != nil {
		h.HandleError(w, r, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify that the required fields are present.
	if rBody.Name == nil || rBody.UnitType == nil || rBody.Quantity == nil {
		hErr := errors.New("request body missing required fields")
		h.HandleError(w, r, hErr, "Invalid request body", http.StatusBadRequest)
		return
	}

	if rBody.Gtin != nil { // if not omitted
		if err := h.checkGtinUniqueness(ctx, *rBody.Gtin, -1); err != nil {
			h.HandleError(w, r, err, "GTIN already exists", http.StatusConflict)
			return
		}
	}

	// Convert request fields to database parameters
	dbParams := convertToCreateFoodEntryParams(rBody, userID)

	// Register the food entry with the database
	foodEntry, err := h.queries.CreateFoodEntry(ctx, dbParams)
	if err != nil {
		h.HandleError(w, r, err, "db error", http.StatusInternalServerError)
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

	// Extract user ID from context - this should be set by the authentication middleware
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
		h.HandleError(w, r, err, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.HandleError(w, r, err, "Couldn't decode request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(rBody); err != nil {
		h.HandleError(w, r, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Consider locking editing to owner of the food entry
	// `if existingFood.CreatedBy != userID`

	if rBody.Gtin != nil { // if not omitted
		if err := h.checkGtinUniqueness(ctx, *rBody.Gtin, -1); err != nil {
			h.HandleError(w, r, err, "GTIN already exists", http.StatusConflict)
			return
		}
	}
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
