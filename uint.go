package into

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

// CanUint returns true if the given value can be coerced to a signed integer.
// [Uint] will succeed without panicking if CanUint returns true.
//
// See: [Uint] for supported types.
func CanUint(x any, options ...Option) bool {
	switch x := x.(type) {
	case uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8, nil:
		return true
	case string, *string, []byte, rune, *rune:
		if !should(options, convertStrings) {
			return false
		}
		if !should(options, checkMarshal) {
			return true
		}
		str := String(x)
		_, err := strconv.ParseUint(str, 10, 64)
		return err == nil
	case encoding.TextMarshaler:
		if !should(options, convertStrings) {
			return false
		}
		if !should(options, checkMarshal) {
			return true
		}
		text, err := x.MarshalText()
		if err != nil {
			return false
		}
		_, err = strconv.ParseUint(string(text), 10, 64)
		return err == nil
	}

	if !should(options, skipReflect) {
		rt := reflect.TypeOf(x)
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		switch rt.Kind() {
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			return true
		case reflect.String:
			return should(options, convertStrings)
		case reflect.Slice:
			if !should(options, convertStrings) {
				return false
			}
			switch rt.Elem().Kind() {
			case reflect.Uint8, reflect.Int32: // []byte, []rune
				return true
			}
			return false
		}
	}

	return false
}

// Uint coerces x into an unsigned integer, supporting the following types:
//   - uint, uint64, uint32, uint16, uint8
//   - *uint, *uint64, *uint32, *uint16, *uint8
//   - types with an underlying unsigned integer value or pointers to such types
//   - given [WithConvertStrings], any string-like type supported by [String]
//   - nil
//
// Uint will panic with ErrInvalid if the value cannot be coerced.
func Uint(x any, options ...Option) uint {
	switch x := x.(type) {
	case uint:
		return x
	case uint64:
		return uint(x)
	case uint32:
		return uint(x)
	case uint16:
		return uint(x)
	case uint8:
		return uint(x)
	case *uint:
		if x == nil {
			goto fallback
		}
		return *x
	case *uint64:
		if x == nil {
			goto fallback
		}
		return uint(*x)
	case *uint32:
		if x == nil {
			goto fallback
		}
		return uint(*x)
	case *uint16:
		if x == nil {
			goto fallback
		}
		return uint(*x)
	case *uint8:
		if x == nil {
			goto fallback
		}
		return uint(*x)
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
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			return uint(rv.Uint())
		}
	}

	if should(options, convertStrings) && CanString(x, options...) {
		// unreachable:
		// var str string
		// if should(options, skipReflect) {
		// 	str = String(x, skipReflect)
		// } else {
		// 	str = String(x)
		// }

		str := String(x)
		if str == "" {
			goto fallback
		}
		if n, err := strconv.ParseUint(str, 10, 64); err == nil {
			return uint(n)
		}
	}

	panic(ErrInvalid{Value: x, Type: "uint"})

fallback:
	var null uint
	for _, opt := range options {
		if opt, ok := opt.(fallbackValue); ok {
			null, ok = opt.x.(uint)
			if !ok {
				panic(fmt.Errorf("invalid fallback value: %v (type of %T), must be %s", opt, opt, "int"))
			}
			break
		}
	}
	return null
}
