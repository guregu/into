// Package into provides convenience functions for creating and dereferencing pointers and converting types.
package into

// Value indirects (dereferences) p if non-nil, otherwise returns the zero value.
func Value[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}

// ValueOr indirects (dereferences) p if non-nil, otherwise returns the fallback value.
func ValueOr[T any](p *T, fallback T) T {
	if p != nil {
		return *p
	}
	return fallback
}

// Ptr returns a pointer to v.
func Ptr[T any](v T) *T {
	return &v
}
