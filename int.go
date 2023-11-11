package into

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

// CanInt returns true if the given value can be coerced to a signed integer.
// [Int] will succeed without panicking if CanInt returns true.
//
// See: [Int] for supported types.
func CanInt(x any, options ...Option) bool {
	switch x := x.(type) {
	case int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8, nil:
		return true
	case string, *string, []byte, []rune, fmt.Stringer:
		if !should(options, convertStrings) {
			return false
		}
		if should(options, skipMarshalCheck) {
			return true
		}
		str := String(x)
		_, err := strconv.Atoi(str)
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
		_, err = strconv.Atoi(string(text))
		return err == nil
	}

	if !should(options, skipReflect) {
		rv := reflect.ValueOf(x)
		for rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			return true
		case reflect.String:
			if !should(options, convertStrings) {
				return false
			}
			if should(options, skipMarshalCheck) {
				return true
			}
			_, err := strconv.Atoi(rv.String())
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
				_, err := strconv.Atoi(string(rv.Bytes()))
				return err == nil
			case reflect.Int32: // []rune
				if should(options, skipMarshalCheck) {
					return true
				}
				_, err := strconv.Atoi(string(rv.Convert(runesType).Interface().([]rune)))
				return err == nil
			}
			return false
		}
	}

	return false
}

// Int coerces x into a signed integer, supporting the following types:
//   - int, int64, int32 (and rune), int16, int8
//   - *int, *int64, *int32 (and *rune), *int16, *int8
//   - types with an underlying signed integer value or pointers to such types
//   - given [WithConvertStrings], any string-like type supported by [String]
//   - nil
//
// Int will panic with ErrInvalid if the value cannot be coerced.
func Int(x any, options ...Option) int {
	switch x := x.(type) {
	case int:
		return x
	case int64:
		return int(x)
	case int32:
		return int(x)
	case int16:
		return int(x)
	case int8:
		return int(x)
	case *int:
		if x == nil {
			goto fallback
		}
		return *x
	case *int64:
		if x == nil {
			goto fallback
		}
		return int(*x)
	case *int32:
		if x == nil {
			goto fallback
		}
		return int(*x)
	case *int16:
		if x == nil {
			goto fallback
		}
		return int(*x)
	case *int8:
		if x == nil {
			goto fallback
		}
		return int(*x)
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
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			return int(rv.Int())
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
		n, err := strconv.Atoi(str)
		if err != nil {
			panic(ErrInvalid{Value: x, Type: "int", Cause: err})
		}
		return n
	}

	panic(ErrInvalid{Value: x, Type: "int"})

fallback:
	var null int
	for _, opt := range options {
		if opt, ok := opt.(fallbackValue); ok {
			null, ok = opt.x.(int)
			if !ok {
				panic(fmt.Errorf("invalid fallback value: %v (type of %T), must be %s", opt, opt, "int"))
			}
			break
		}
	}
	return null
}
