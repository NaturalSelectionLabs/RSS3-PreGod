package util

func TrimQuote(s string) string {
	return s[1 : len(s)-1]
}
