package testdata

import "database/sql"

func Uint64Pointer(i uint64) *uint64 {
	return &i
}

func Uint32Pointer(i uint32) *uint32 {
	return &i
}

func Float32Pointer(f float32) *float32 {
	return &f
}

func Float64Pointer(f float64) *float64 {
	return &f
}

func StringPointer(s string) *string {
	return &s
}

func Int32Pointer(i int32) *int32 {
	return &i
}

func NullInt64Uint32Ptr(u *uint32) sql.NullInt64 {
	ret := sql.NullInt64{Int64: 0, Valid: false}
	if u != nil {
		ret.Valid = true
		ret.Int64 = int64(*u)
	}
	return ret
}

func NullInt64Int32Ptr(u *int32) sql.NullInt64 {
	ret := sql.NullInt64{Int64: 0, Valid: false}
	if u != nil {
		ret.Valid = true
		ret.Int64 = int64(*u)
	}
	return ret
}

func NullStringStringPtr(s *string) sql.NullString {
	ret := sql.NullString{String: "", Valid: false}
	if s != nil {
		ret.Valid = true
		ret.String = *s
	}
	return ret
}

func NullFloat64Float32Ptr(f *float32) sql.NullFloat64 {
	ret := sql.NullFloat64{Float64: 0, Valid: false}
	if f != nil {
		ret.Valid = true
		ret.Float64 = float64(*f)
	}
	return ret
}

func Uint32PtrNullInt64(i sql.NullInt64) *uint32 {
	if !i.Valid {
		return nil
	}
	return Uint32Pointer(uint32(i.Int64))
}

func Int32PtNullInt64(i sql.NullInt64) *int32 {
	if !i.Valid {
		return nil
	}
	return Int32Pointer(int32(i.Int64))
}

func StringPtrNullString(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return StringPointer(s.String)
}

func Float32PtrNullFloat(f sql.NullFloat64) *float32 {
	if !f.Valid {
		return nil
	}
	return Float32Pointer(float32(f.Float64))
}
