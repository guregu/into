package into

import "slices"

// Option is a configuration parameter.
// Use the With... functions to specify options.
type Option interface{ isOption() }

func should(options []Option, opt Option) bool {
	return slices.Contains[[]Option, Option](options, opt)
}

type flags int

func (flags) isOption() {}

const (
	convertStrings flags = iota
	skipReflect
	checkMarshal
)

type fallbackValue struct{ x any }

func (fallbackValue) isOption() {}

// WithFallback specifies a fallback value when coercing nil input.
// By default, the zero value is returned.
func WithFallback(fallback any) Option {
	return fallbackValue{fallback}
}

// WithConvertStrings enables conversion of strings during type coercion.
func WithConvertStrings() Option {
	return convertStrings
}

// WithoutReflection will skip using reflection to coerce values.
// Using this disables support for nonstandard types (e.g. custom subtypes of int or string).
func WithoutReflection() Option {
	return skipReflect
}

// WithMarshalerCheck is an option that enables an additional check in [CanString] and similar,
// which runs [encoding.TextMarshaler]'s marshal or [strconv] conversions to ensure they don't return errors.
func WithMarshalerCheck() Option {
	return checkMarshal
}
