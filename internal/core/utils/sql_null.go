package utils

import (
	"database/sql"
	"time"
)

// ToNullString converts a string to sql.NullString.
// If the string is empty, it returns a null representation.
func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// ToNullInt32 converts an int32 to sql.NullInt32.
// If the value is 0, it is considered null. (Adjust based on business logic)
func ToNullInt32(i int32) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: i, Valid: true}
}

// ToNullInt64 converts an int64 to sql.NullInt64.
// If the value is 0, it is considered null.
func ToNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}

// ToNullFloat64 converts a float64 to sql.NullFloat64.
// If the value is 0, it is considered null.
func ToNullFloat64(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

// ToNullBool converts a bool pointer to sql.NullBool.
// We use pointer here because bool zero value is false, which might be a valid non-null value.
func ToNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

// ToNullTime converts a time.Time to sql.NullTime.
// If the time is zero (0001-01-01), it returns a null representation.
func ToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}
