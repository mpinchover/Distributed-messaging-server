package utils

func ToStrPtr(s string) *string {
	return &s
}

func ToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
