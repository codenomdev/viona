package util

import "strconv"

func ParseStringInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ParseInt64String(d int64) string {
	return strconv.FormatInt(d, 10)
}

func ParseUint64String(d uint64) string {
	return strconv.FormatUint(d, 10)
}
