package typeconvert

import "github.com/jackc/pgx/v5/pgtype"

func PtrStringToPgtypeText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func PgtypeTextToString(t pgtype.Text) string { // FIXME: consider ptr
	if !t.Valid {
		return ""
	}
	return t.String
}

func PtrInt32ToPgtypeInt4(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *i, Valid: true}
}

func PgtypeInt4ToInt(i pgtype.Int4) int32 { // FIXME: consider ptr
	if !i.Valid {
		return 0
	}
	return i.Int32
}

func PtrFloat32ToPgtypeFloat4(f *float32) pgtype.Float4 {
	if f == nil {
		return pgtype.Float4{Valid: false}
	}
	return pgtype.Float4{Float32: *f, Valid: true}
}

func PgtypeFloat4ToFloat32(pf pgtype.Float4) float32 { // FIXME: consider ptr
	if !pf.Valid {
		return 0.0
	}
	return pf.Float32
}

func PtrBoolToPgtypeBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func PgtypeBoolToBool(b pgtype.Bool) bool {
	// FIXME: consider ptr
	if !b.Valid {
		return false
	}
	return b.Bool
}
