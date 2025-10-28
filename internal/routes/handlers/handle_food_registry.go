package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/PlateOps/util/ctxutils"
	"github.com/Laelapa/PlateOps/util/parse"
	"github.com/Laelapa/PlateOps/util/typeconvert"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// food represents the structure of the request body for creating or updating a food entry.
// All fields are set to be optional to accommodate the PATCH method.
// For fields that are necessary for other methods verify the request body in the handler.
// Check for presence of a field in a request by checking if the pointer is nil.
type food struct {
	ID                     *int32   `json:"id,omitempty"` // ID is irrelevant for POST and PATCH but is returned by GET methods.
	Name                   *string  `json:"name"                     validate:"omitempty,max=255"`
	Gtin                   *string  `json:"gtin"                     validate:"omitempty,numeric,max=14"` // GTIN-14
	Category               *string  `json:"category"                 validate:"omitempty,max=255"`
	Description            *string  `json:"description"              validate:"omitempty,max=5000"`
	UnitType               *string  `json:"unit_type"                validate:"omitempty,oneof=grams ml items portions"`
	Quantity               *int32   `json:"quantity"`
	PortionCount           *int32   `json:"portion_count"            validate:"omitempty,min=1"` // Number of portions in the food item
	ExpirationAfterOpening *int32   `json:"expiration_after_opening" validate:"omitempty,min=0"` // in days
	NutrientsPerItem       *bool    `json:"nutrients_per_item"`                                  // true: per item/portion, false: per 100g
	Calories               *float32 `json:"calories"                 validate:"omitempty,min=0"`
	Fats                   *float32 `json:"fats"                     validate:"omitempty,min=0"`
	Saturated              *float32 `json:"saturated"                validate:"omitempty,min=0"`
	Carbs                  *float32 `json:"carbs"                    validate:"omitempty,min=0"`
	Sugars                 *float32 `json:"sugars"                   validate:"omitempty,min=0"`
	Protein                *float32 `json:"protein"                  validate:"omitempty,min=0"`
	Fiber                  *float32 `json:"fiber"                    validate:"omitempty,min=0"`
	Sodium                 *float32 `json:"sodium"                   validate:"omitempty,min=0"`
}

type responseFood struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	ProductID int32  `json:"product_id,omitempty"`
}

var ErrGtinAlreadyExists = errors.New("GTIN already exists for another food entry")
var ErrQueryFailed = errors.New("failed to execute database query")

// TODO: Refactor GET handlers to reduce code duplication
func (h *Handler) HandleGetFoodByID(w http.ResponseWriter, r *http.Request) {

	h.logger.LogRequestInfo("Fetching food entry by ID", r)

	ctx := r.Context()
	productID, err := parse.ID(r.PathValue("id"))
	if err != nil {
		h.HandleError(w, r, err, "Invalid product ID", http.StatusBadRequest)
		return
	}

	foodEntry, err := h.queries.GetFoodEntryById(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.HandleError(w, r, err, "Food entry not found", http.StatusNotFound)
			return
		}
		h.HandleError(w, r, err, "db error", http.StatusInternalServerError)
		return
	}

	resp := convertToFoodResponse(foodEntry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.LogAppError("Failed to encode response", err)
	}
}

func (h *Handler) HandleGetFoodByGtin(w http.ResponseWriter, r *http.Request) {

	h.logger.LogRequestInfo("Fetching food entry by GTIN", r)

	ctx := r.Context()
	productGTIN, err := parse.GTIN(r.PathValue("gtin"))
	if err != nil {
		h.HandleError(w, r, err, "Invalid product GTIN", http.StatusBadRequest)
		return
	}

	foodEntry, err := h.queries.GetFoodEntryByGtin(ctx, productGTIN)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.HandleError(w, r, err, "Food entry not found", http.StatusNotFound)
			return
		}
		h.HandleError(w, r, err, "db error", http.StatusInternalServerError)
		return
	}

	resp := convertToFoodResponse(foodEntry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.LogAppError("Failed to encode response", err)
	}
}

func (h *Handler) HandleGetFoodsByNameContains(w http.ResponseWriter, r *http.Request) {

	h.logger.LogRequestInfo("Fetching food entries by name", r)

	ctx := r.Context()
	foodName := r.PathValue("name")

	// Validate the name parameter
	if foodName == "" {
		h.HandleError(w, r, errors.New("missing food name"), "Invalid food name", http.StatusBadRequest)
		return
	}

	// Optional: Add length validation to match your struct validation
	if len(foodName) > 255 {
		h.HandleError(w, r, errors.New("food name too long"), "Invalid food name", http.StatusBadRequest)
		return
	}

	// Convert string to pgtype.Text for database query
	nameParam := typeconvert.StringToPgtypeText(foodName)

	foodEntries, err := h.queries.GetFoodEntriesByNameContains(ctx, nameParam)
	if err != nil {
		h.HandleError(w, r, err, "db error", http.StatusInternalServerError)
		return
	}

	// Convert database results to response format
	var foods []food
	for _, entry := range foodEntries {
		foods = append(foods, convertToFoodResponse(entry))
	}

	// Create response structure
	response := struct {
		Foods []food `json:"foods"`
		Count int    `json:"count"`
	}{
		Foods: foods,
		Count: len(foods),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.LogAppError("Failed to encode response", err)
	}

	h.logger.LogAppInfo(
		"Food entries retrieved successfully",
		zap.String("name", foodName),
		zap.Int("count", len(foods)),
	)
}

func (h *Handler) HandlePostFood(w http.ResponseWriter, r *http.Request) {

	var rBody food
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

	rBody.ID = nil // ID is not set for POST requests, it will be generated by the database.

	if rBody.Gtin != nil { // if not omitted
		if err := h.checkGtinUniqueness(ctx, rBody.Gtin, -1); err != nil {
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

	var rBody food
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

	rBody.ID = nil // ID is not set for PATCH requests, it is determined by the URL.

	// TODO: Consider locking editing to owner of the food entry
	// `if existingFood.CreatedBy != userID`

	if rBody.Gtin != nil { // if not omitted
		if err := h.checkGtinUniqueness(ctx, rBody.Gtin, -1); err != nil {
			h.HandleError(w, r, err, "GTIN already exists", http.StatusConflict)
			return
		}
	}

	// Convert request fields to database parameters
	dbParams := convertToUpdateFoodEntryParams(rBody, userID, productID)

	dbErr := h.queries.UpdateFoodEntry(ctx, dbParams)
	if dbErr != nil {
		if errors.Is(dbErr, pgx.ErrNoRows) {
			h.HandleError(w, r, dbErr, "Food entry not found", http.StatusNotFound)
			return
		}
		h.HandleError(w, r, dbErr, "db error", http.StatusInternalServerError)
		return
	}

	h.publishFoodPatchEvent(context.Background(), productID)

	resp := responseFood{
		Success:   true,
		Message:   "Food entry updated successfully",
		ProductID: productID,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.LogAppError("Failed to encode response", err)
	}

	h.logger.LogAppInfo(
		"Food entry updated successfully",
		zap.Int32("product_id", productID),
		zap.String("user_id", typeconvert.PgtypeUUIDToString(userID)),
	)
}
