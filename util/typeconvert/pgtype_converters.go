package typeconvert

import "github.com/jackc/pgx/v5/pgtype"

func StringToPgtypeText(s string) pgtype.Text {
    return pgtype.Text{String: s, Valid: s != ""}
}

func PtrStringToPgtypeText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func PgtypeTextToPtrString(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func PtrInt32ToPgtypeInt4(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *i, Valid: true}
}

func PgtypeInt4ToPtrInt32(i pgtype.Int4) *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

func PtrFloat32ToPgtypeFloat4(f *float32) pgtype.Float4 {
	if f == nil {
		return pgtype.Float4{Valid: false}
	}
	return pgtype.Float4{Float32: *f, Valid: true}
}

func PgtypeFloat4ToPtrFloat32(pf pgtype.Float4) *float32 {
	if !pf.Valid {
		return nil
	}
	return &pf.Float32
}

func PtrBoolToPgtypeBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func PgtypeBoolToPtrBool(b pgtype.Bool) *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}
