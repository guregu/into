package into

import "fmt"

// ErrInvalid is an error returned when into could not convert a given value.
type ErrInvalid struct {
	Value any
	Type  string

	// Cause is non-nil if a fallible conversion returns an error (such as string conversion or TextMarshaler).
	Cause error
}

func (err ErrInvalid) Error() string {
	var extra string
	if err.Cause != nil {
		extra = "; " + err.Cause.Error()
	}
	return fmt.Sprintf("into: value %v of type %T is not a %s%s", err.Value, err.Value, err.Type, extra)
}

func (err ErrInvalid) Unwrap() error {
	return err.Cause
}
