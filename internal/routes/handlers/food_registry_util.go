package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/util"
	"github.com/Laelapa/PlateOps/util/typeconvert"
)

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
		Topic: "food.updates",
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
