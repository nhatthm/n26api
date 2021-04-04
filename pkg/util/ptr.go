package util

// StringPtr returns the a pointer of string.
func StringPtr(str string) *string {
	return &str
}

// Int64Ptr returns the a pointer of int64.
func Int64Ptr(i int64) *int64 {
	return &i
}
