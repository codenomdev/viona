package util

import (
	"time"

	"github.com/gookit/goutil/strutil"
)

func TimeNow() time.Time {
	return time.Now().In(time.FixedZone("WIB", 7*60*60))
}

// Generate random token.
// Len 32 characters include alphanum, -, _
func GenerateRandomToken() (string, error) {
	return strutil.RandomString(32)
}

// Generate random string Tpl
// len n characters
func GenerateStrTpl(n int, tpl string) string {
	if tpl == "" {
		tpl = strutil.AlphaNum3
	}
	return strutil.RandWithTpl(n, tpl)
}

// Generate time unix with add 2 hours to Mili.
func GenerateTimeHours() time.Time {
	return TimeNow().Add(2 * time.Hour)
}

// GenerateTimeOneMonth()
func GenerateTimeOneMonth() time.Time {
	return TimeNow().Add(30 * 24 * time.Hour)
}

// Generate time monthly
func GenerateTimeMonthly() time.Time {
	return TimeNow().AddDate(0, 0, 29)
}

// Generate time yearly
func GenerateTimeYearly() time.Time {
	return TimeNow().AddDate(0, 11, 29)
}
