package typeconvert

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// GoogleUUIDToPgtypeUUID converts a Google UUID to a pgtype.UUID.
// It returns a pgtype.UUID with the Valid field set to true.
// If the Google UUID is nil, it sets the Valid field to false.
func GoogleUUIDToPgtypeUUID(gUUID uuid.UUID) pgtype.UUID {
	if gUUID == uuid.Nil {
		return pgtype.UUID{
			Valid: false,
		}		
	}

	return pgtype.UUID{
		Bytes: gUUID,
		Valid: true,
	}
}

// PgtypeUUIDToGoogleUUID converts a pgtype.UUID to a Google UUID.
// If the pgtype.UUID is invalid, it returns uuid.Nil.
func PgtypeUUIDToGoogleUUID(pgUUID pgtype.UUID) uuid.UUID {
	if !pgUUID.Valid {
		return uuid.Nil
	}

	return pgUUID.Bytes
}

// PgtypeUUIDToString converts a pgtype.UUID to a string.
// If the pgtype.UUID is invalid, it returns an empty string.
func PgtypeUUIDToString(pgUUID pgtype.UUID) string {

	if gUUID := PgtypeUUIDToGoogleUUID(pgUUID); gUUID == uuid.Nil {
		return ""
	} else {
		return gUUID.String()
	}
}

// StringToPgtypeUUID converts a string to a pgtype.UUID.
// If the string is empty or cannot be parsed as a UUID, it returns a pgtype.UUID with the field `Valid: false`.
func StringToPgtypeUUID(s string) pgtype.UUID {
	if s == "" {
		return pgtype.UUID{
			Valid: false,
		}
	}

	gUUID, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{
			Valid: false,
		}
	}

	return GoogleUUIDToPgtypeUUID(gUUID)
}