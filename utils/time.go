package utils

import "time"

func Format(unix int64) string {
	t := time.Unix(unix, 0)

	return t.Format("2006-01-02")
}

func Month(unix int64) string {
	t := time.Unix(unix, 0)

	return t.Format("1")
}

func Year(unix int64) string {
	t := time.Unix(unix, 0)

	return t.Format("2006")
}

func CMonth(unix int64) string {
	t := time.Unix(unix, 0)

	return t.Format("01-02")
}

func Str2Unix(layout, tstr string) int64 {
	tm, _ := time.Parse(layout, tstr)
	return tm.Unix()
}
