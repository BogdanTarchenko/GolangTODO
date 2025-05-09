package utils

// Ptr returns a pointer to the given value.
// This is a generic function that can be used with any type.
func Ptr[T any](v T) *T {
	return &v
} 