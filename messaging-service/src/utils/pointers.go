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

func ToInt64Ptr(i int64) *int64 {
	return &i
}

func ToInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func ToFloat64Ptr(i float64) *float64 {
	return &i
}

func ToFloat64(i *float64) float64 {
	if i == nil {
		return 0
	}
	return *i
}
