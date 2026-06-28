package util

import (
	"database/sql"
	"strings"
	"time"
)

// StringToNullInt64 converts a string to sql.NullInt64.
// If the string is empty or unparseable, it returns NullInt64 with Valid = false.
func StringToNullInt64(s string) (sql.NullInt64, error) {
	// 1. Check whether the string is empty or just spaces (NULL representation)
	if strings.TrimSpace(s) == "" {
		return sql.NullInt64{Valid: false}, nil
	}

	// 2. Parse string to int64
	num, err := ParseStringInt64(s)

	if err != nil {
		// If parsing fails (e.g., the string contains "abc")
		// we still return NullInt64 with Valid=false,
		// but also return an error so the calling logic knows there was a parsing problem.
		return sql.NullInt64{Valid: false}, err
	}

	// 3. Conversion successful: Create struct sql.NullInt64
	return sql.NullInt64{
		Int64: num,
		Valid: true, // Set Valid = true because the value was found and is valid.
	}, nil
}

func ToNullInt64(v int64) sql.NullInt64 {
	if v == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: v, Valid: true}
}

func TimeToNull(times time.Time) sql.NullTime {
	if times.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: times, Valid: true}
}
