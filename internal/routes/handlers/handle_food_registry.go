package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/util"
	"github.com/Laelapa/PlateOps/util/ctxutils"
	"github.com/Laelapa/PlateOps/util/parse"
	"github.com/Laelapa/PlateOps/util/typeconvert"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	h.publishFoodPatchEvent(ctx, productID)

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

func convertToFoodResponse(foodEntry repository.FoodRegistry) food {

	return food{
		ID:                     &foodEntry.ProductID,
		Name:                   typeconvert.PgtypeTextToPtrString(foodEntry.Name),
		Gtin:                   typeconvert.PgtypeTextToPtrString(foodEntry.Gtin),
		Category:               typeconvert.PgtypeTextToPtrString(foodEntry.Category),
		Description:            typeconvert.PgtypeTextToPtrString(foodEntry.Description),
		UnitType:               typeconvert.PgtypeTextToPtrString(foodEntry.UnitType),
		Quantity:               typeconvert.PgtypeInt4ToPtrInt32(foodEntry.Quantity),
		PortionCount:           typeconvert.PgtypeInt4ToPtrInt32(foodEntry.PortionCount),
		ExpirationAfterOpening: typeconvert.PgtypeInt4ToPtrInt32(foodEntry.ExpirationAfterOpening),
		NutrientsPerItem:       typeconvert.PgtypeBoolToPtrBool(foodEntry.NutrientsPerItem),
		Calories:               typeconvert.PgtypeFloat4ToPtrFloat32(foodEntry.Calories),
		Fats:                   typeconvert.PgtypeFloat4ToPtrFloat32(foodEntry.Fats),
		Saturated:              typeconvert.PgtypeFloat4ToPtrFloat32(foodEntry.Saturated),
		Carbs:                  typeconvert.PgtypeFloat4ToPtrFloat32(foodEntry.Carbs),
		Sugars:                 typeconvert.PgtypeFloat4ToPtrFloat32(foodEntry.Sugars),
		Protein:                typeconvert.PgtypeFloat4ToPtrFloat32(foodEntry.Protein),
		Fiber:                  typeconvert.PgtypeFloat4ToPtrFloat32(foodEntry.Fiber),
		Sodium:                 typeconvert.PgtypeFloat4ToPtrFloat32(foodEntry.Sodium),
		// Note: metadata are not included in the response.
	}
}

func convertToCreateFoodEntryParams(
	req food,
	userID pgtype.UUID,
) repository.CreateFoodEntryParams {

	return repository.CreateFoodEntryParams{
		Name:                   typeconvert.PtrStringToPgtypeText(req.Name),
		Gtin:                   typeconvert.PtrStringToPgtypeText(req.Gtin),
		Category:               typeconvert.PtrStringToPgtypeText(req.Category),
		Description:            typeconvert.PtrStringToPgtypeText(req.Description),
		UnitType:               typeconvert.PtrStringToPgtypeText(req.UnitType),
		Quantity:               typeconvert.PtrInt32ToPgtypeInt4(req.Quantity),
		PortionCount:           typeconvert.PtrInt32ToPgtypeInt4(req.PortionCount),
		ExpirationAfterOpening: typeconvert.PtrInt32ToPgtypeInt4(req.ExpirationAfterOpening),
		NutrientsPerItem:       typeconvert.PtrBoolToPgtypeBool(req.NutrientsPerItem),
		Calories:               typeconvert.PtrFloat32ToPgtypeFloat4(req.Calories),
		Fats:                   typeconvert.PtrFloat32ToPgtypeFloat4(req.Fats),
		Saturated:              typeconvert.PtrFloat32ToPgtypeFloat4(req.Saturated),
		Carbs:                  typeconvert.PtrFloat32ToPgtypeFloat4(req.Carbs),
		Sugars:                 typeconvert.PtrFloat32ToPgtypeFloat4(req.Sugars),
		Protein:                typeconvert.PtrFloat32ToPgtypeFloat4(req.Protein),
		Fiber:                  typeconvert.PtrFloat32ToPgtypeFloat4(req.Fiber),
		Sodium:                 typeconvert.PtrFloat32ToPgtypeFloat4(req.Sodium),
		CreatedBy:              userID,
	}
}

func convertToUpdateFoodEntryParams(
	req food,
	userID pgtype.UUID,
	productID int32,
) repository.UpdateFoodEntryParams {

	return repository.UpdateFoodEntryParams{
		Name:                   typeconvert.PtrStringToPgtypeText(req.Name),
		Gtin:                   typeconvert.PtrStringToPgtypeText(req.Gtin),
		Category:               typeconvert.PtrStringToPgtypeText(req.Category),
		Description:            typeconvert.PtrStringToPgtypeText(req.Description),
		UnitType:               typeconvert.PtrStringToPgtypeText(req.UnitType),
		Quantity:               typeconvert.PtrInt32ToPgtypeInt4(req.Quantity),
		PortionCount:           typeconvert.PtrInt32ToPgtypeInt4(req.PortionCount),
		ExpirationAfterOpening: typeconvert.PtrInt32ToPgtypeInt4(req.ExpirationAfterOpening),
		NutrientsPerItem:       typeconvert.PtrBoolToPgtypeBool(req.NutrientsPerItem),
		Calories:               typeconvert.PtrFloat32ToPgtypeFloat4(req.Calories),
		Fats:                   typeconvert.PtrFloat32ToPgtypeFloat4(req.Fats),
		Saturated:              typeconvert.PtrFloat32ToPgtypeFloat4(req.Saturated),
		Carbs:                  typeconvert.PtrFloat32ToPgtypeFloat4(req.Carbs),
		Sugars:                 typeconvert.PtrFloat32ToPgtypeFloat4(req.Sugars),
		Protein:                typeconvert.PtrFloat32ToPgtypeFloat4(req.Protein),
		Fiber:                  typeconvert.PtrFloat32ToPgtypeFloat4(req.Fiber),
		Sodium:                 typeconvert.PtrFloat32ToPgtypeFloat4(req.Sodium),
		ProductID:              productID,
		UpdatedBy:              userID,
	}
}

func (h *Handler) checkGtinUniqueness(ctx context.Context, gtin *string, productID int32) error {

	// Check if GTIN is already used by another food entry
	existingFood, err := h.queries.GetFoodEntryByGtin(
		ctx,
		typeconvert.PtrStringToPgtypeText(gtin),
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

func (h *Handler) publishFoodPatchEvent(ctx context.Context, productID int32) {
	if h.kafkaClient == nil {
		h.logger.LogAppWarn("Kafka client is not initialized, skipping event publishing")
		return
	}

	eventData := map[string]interface{}{
		"event":   "food.updated",
		"food_id": productID,
	}

	eventBytes, err := json.Marshal(eventData)
	if err != nil {
		h.logger.LogAppError("Failed to marshal food update event", err)
		return
	}

	record := &kgo.Record{
		Topic: "food_updates",
		Value: eventBytes,
	}

	h.kafkaClient.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			h.logger.LogAppError("Failed to send food update event", err)
		} else {
			h.logger.LogAppInfo("Food update event published successfully",
				zap.Int32("food_id", productID),
				zap.String("topic", r.Topic),
				zap.Int32("partition", r.Partition),
				zap.Int64("offset", r.Offset),
				zap.Time("kafka_timestamp", r.Timestamp))
		}
	})
}
