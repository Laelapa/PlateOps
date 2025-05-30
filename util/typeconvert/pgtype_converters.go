package typeconvert

import "github.com/jackc/pgx/v5/pgtype"

func StringToPgtypeText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func PgtypeTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func IntToPgtypeInt4(i int) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(i), Valid: true}
}

func PgtypeInt4ToInt(i pgtype.Int4) int {
	if !i.Valid {
		return 0
	}
	return int(i.Int32)
}

func Float32ToPgtypeFloat4(f float32) pgtype.Float4 {
	return pgtype.Float4{Float32: f, Valid: true}
}

func PgtypeFloat4ToFloat32(pf pgtype.Float4) float32 {
	if !pf.Valid {
		return 0.0
	}
	return pf.Float32
}

func BoolToPgtypeBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func PgtypeBoolToBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}
