package into

import (
	"encoding"
	"fmt"
	"reflect"
)

// CanString returns true if the given value can be coerced to a signed integer.
// [String] will succeed without panicking if CanString returns true.
//
// See: [String] for supported types.
func CanString(x any, options ...Option) bool {
	switch x := x.(type) {
	case string, *string, []byte, rune, *rune, []rune, fmt.Stringer, nil:
		return true
	case encoding.TextMarshaler:
		if !should(options, checkMarshal) {
			return true
		}
		_, err := x.MarshalText()
		return err == nil
	}

	if !should(options, skipReflect) {
		rt := reflect.TypeOf(x)
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		switch rt.Kind() {
		case reflect.String:
			return true
		case reflect.Slice:
			switch rt.Elem().Kind() {
			case reflect.Uint8, reflect.Int32: // []byte, []rune
				return true
			}
		}
		return false
	}

	return false
}

// String coerces x into a string, supporting the following types:
//   - string, []byte, rune, []rune
//   - *string, *rune
//   - types with an underlying value of string, []byte, rune, or []rune, unless [WithoutReflection] is used
//   - [encoding.TextMarshaler]
//   - [fmt.Stringer]
//   - nil
//
// String will panic with ErrInvalid if the value cannot be coerced or TextMarshaler fails.
func String(x any, options ...Option) string {
	switch x := x.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case rune:
		return string(x)
	case []rune:
		return string(x)
	case *string:
		if x == nil {
			goto fallback
		}
		return *x
	case *rune:
		if x == nil {
			goto fallback
		}
		return string(*x)
	case encoding.TextMarshaler:
		bs, err := x.MarshalText()
		if err != nil {
			panic(ErrInvalid{Value: x, Type: "string", Cause: err})
		}
		return string(bs)
	case fmt.Stringer:
		return x.String()
	case nil:
		goto fallback
	}

	if !should(options, skipReflect) {
		rv := reflect.ValueOf(x)
		for rv.Kind() == reflect.Pointer {
			if rv.IsNil() {
				goto fallback
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.String:
			return rv.String()
		case reflect.Slice:
			switch rv.Type().Elem().Kind() {
			case reflect.Uint8:
				if rv.IsNil() {
					goto fallback
				}
				return string(rv.Bytes())
			case reflect.Int32:
				if rv.IsNil() {
					goto fallback
				}
				return string(rv.Convert(runesType).Interface().([]rune))
			default:
				panic(ErrInvalid{Value: x, Type: "string"})
			}
		}
	}

	panic(ErrInvalid{Value: x, Type: "string"})

fallback:
	var null string
	for _, opt := range options {
		if opt, ok := opt.(fallbackValue); ok {
			null, ok = opt.x.(string)
			if !ok {
				panic(fmt.Errorf("invalid fallback value: %v (type of %T), must be %s", opt, opt, "string"))
			}
			break
		}
	}
	return null
}

var runesType = reflect.TypeOf([]rune{})
