package typeconvert

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TimeToPgtypeTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:             t,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

// TODO: Add a function to convert pgtype.Timestamp to time.Time,
// have to consider the InfinityModifier.
