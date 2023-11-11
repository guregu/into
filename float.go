package into

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

// CanFloat returns true if the given value can be coerced to a float.
// See: [Float] for supported types.
func CanFloat(x any, options ...Option) bool {
	switch x := x.(type) {
	case float64, float32, *float64, *float32, nil:
		return true
	case string, *string, []byte, []rune, fmt.Stringer:
		if !should(options, convertStrings) {
			return false
		}
		if should(options, skipMarshalCheck) {
			return true
		}
		str := String(x)
		_, err := strconv.ParseFloat(str, 64)
		return err == nil
	case encoding.TextMarshaler:
		if !should(options, convertStrings) {
			return false
		}
		if should(options, skipMarshalCheck) {
			return true
		}
		text, err := x.MarshalText()
		if err != nil {
			return false
		}
		_, err = strconv.ParseFloat(string(text), 64)
		return err == nil
	}

	if !should(options, skipReflect) {
		rv := reflect.ValueOf(x)
		for rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Float64, reflect.Float32:
			return true
		case reflect.String:
			if !should(options, convertStrings) {
				return false
			}
			if should(options, skipMarshalCheck) {
				return true
			}
			_, err := strconv.ParseFloat(rv.String(), 64)
			return err == nil
		case reflect.Slice:
			if !should(options, convertStrings) {
				return false
			}
			switch rv.Type().Elem().Kind() {
			case reflect.Uint8: // []byte
				if should(options, skipMarshalCheck) {
					return true
				}
				_, err := strconv.ParseFloat(string(rv.Bytes()), 64)
				return err == nil
			case reflect.Int32: // []rune
				if should(options, skipMarshalCheck) {
					return true
				}
				_, err := strconv.ParseFloat(string(rv.Convert(runesType).Interface().([]rune)), 64)
				return err == nil
			}
			return false
		}
	}

	return false
}

// Float coerces x into a float, supporting the following types:
//   - float64, float32
//   - *float64, *float32
//   - types with an underlying float value or pointers to such types
//   - given [WithConvertStrings], any string-like type supported by [String]
//   - nil
//
// Float will panic with ErrInvalid if the value cannot be coerced.
func Float(x any, options ...Option) float64 {
	switch x := x.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case *float64:
		if x == nil {
			goto fallback
		}
		return *x
	case *float32:
		if x == nil {
			goto fallback
		}
		return float64(*x)
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
		case reflect.Float64, reflect.Float32:
			return rv.Float()
		}
	}

	if should(options, convertStrings) && CanString(x) {
		var str string
		if should(options, skipReflect) {
			str = String(x, skipReflect)
		} else {
			str = String(x)
		}
		if str == "" {
			goto fallback
		}
		n, err := strconv.ParseFloat(str, 64)
		if err != nil {
			panic(ErrInvalid{Value: x, Type: "float", Cause: err})
		}
		return n
	}

	panic(ErrInvalid{Value: x, Type: "float"})

fallback:
	var null float64
	for _, opt := range options {
		if opt, ok := opt.(fallbackValue); ok {
			null, ok = opt.x.(float64)
			if !ok {
				panic(fmt.Errorf("invalid fallback value: %v (type of %T), must be %s", opt, opt, "float64"))
			}
			break
		}
	}
	return null
}
