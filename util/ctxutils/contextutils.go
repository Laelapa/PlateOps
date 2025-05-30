package ctxutils

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type contextKey string

const userIDKey contextKey = "userID"

// GetUserIDFromContext is a helper function to extract userID from context in a type-safe manner.
// It returns the userID as a pgtype.UUID and a boolean indicating whether it was found.
func GetUserIDFromContext(ctx context.Context) (pgtype.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(pgtype.UUID)
	return userID, ok
}

// SetUserIDInContext stores a user ID in the context using a type-safe key.
func SetUserIDInContext(ctx context.Context, userID pgtype.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
