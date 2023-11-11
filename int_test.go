package into_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/guregu/into"
)

type myInt int

func TestInt(t *testing.T) {
	t.Parallel()
	tests := table[int]{
		{
			name:  "int",
			input: int(42),
			want:  42,
		},
		{
			name:  "int64",
			input: int64(42),
			want:  42,
		},
		{
			name:  "int32",
			input: int32(42),
			want:  42,
		},
		{
			name:  "int16",
			input: int16(42),
			want:  42,
		},
		{
			name:  "int8",
			input: int8(42),
			want:  42,
		},
		{
			name:  "*int",
			input: into.Ptr(42),
			want:  42,
		},
		{
			name:  "*int64",
			input: into.Ptr(int64(42)),
			want:  42,
		},
		{
			name:  "*int32",
			input: into.Ptr(int32(42)),
			want:  42,
		},
		{
			name:  "*int16",
			input: into.Ptr(int16(42)),
			want:  42,
		},
		{
			name:  "*int8",
			input: into.Ptr(int8(42)),
			want:  42,
		},
		{
			name:  "*int(nil)",
			input: (*int)(nil),
			want:  0,
		},
		{
			name:  "*int64(nil)",
			input: (*int64)(nil),
			want:  0,
		},
		{
			name:  "*int32(nil)",
			input: (*int32)(nil),
			want:  0,
		},
		{
			name:  "*int16(nil)",
			input: (*int16)(nil),
			want:  0,
		},
		{
			name:  "*int8(nil)",
			input: (*int8)(nil),
			want:  0,
		},
		{
			name:  "subtype",
			input: myInt(42),
			want:  42,
		},
		{
			name:  "subtype pointer",
			input: into.Ptr(myInt(42)),
			want:  42,
		},
		{
			name:  "nil subtype pointer",
			input: (*myInt)(nil),
			want:  0,
		},
		{
			name:  "subtype without reflection",
			input: myInt(42),
			want:  0,
			opts:  []into.Option{into.WithoutReflection()},
			err:   into.ErrInvalid{Value: 42, Type: "int"},
		},
		{
			name:  "fallback",
			input: nil,
			want:  42,
			opts:  []into.Option{into.WithFallback(42)},
		},
		{
			name:  "default fallback",
			input: nil,
			want:  0,
		},
		{
			name:  "string conversion",
			input: "42",
			want:  42,
			opts:  []into.Option{into.WithConvertStrings()},
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
			opts:  []into.Option{into.WithConvertStrings()},
		},
		{
			name:  "bad string",
			input: "abc",
			want:  0,
			opts:  []into.Option{into.WithConvertStrings()},
			err:   into.ErrInvalid{Value: "abc", Type: "int", Cause: &strconv.NumError{Func: "Atoi", Num: "abc", Err: strconv.ErrSyntax}},
		},
		{
			name:  "string conversion using subtype",
			input: myString("42"),
			want:  42,
			opts:  []into.Option{into.WithConvertStrings()},
		},
		{
			name:  "string conversion using subtype without reflection",
			input: myString("42"),
			want:  0,
			opts:  []into.Option{into.WithConvertStrings(), into.WithoutReflection()},
			err:   into.ErrInvalid{Value: myString("42"), Type: "int"},
		},
		{
			name:  "invalid type",
			input: struct{}{},
			want:  0,
			opts:  []into.Option{into.WithoutReflection()},
			err:   into.ErrInvalid{Value: struct{}{}, Type: "int"},
		},
	}
	tests.Run(t, into.Int)
}

func TestCanInt(t *testing.T) {
	t.Parallel()

	good := []any{
		int(42),
		int64(42),
		int32(42),
		int16(42),
		int8(42),
		myInt(42),
		new(int),
		new(myInt),
		nil,
	}
	for _, v := range good {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if !into.CanInt(v) {
				t.Error("failed but should have succeeded for value:", v)
			}
		})
	}

	bad := []any{
		float64(1),
		float32(1),
		complex(1, 1),
		struct{}{},
		my42Marshaler{},
		string("42"),
		[]int64{1},
	}
	for _, v := range bad {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if into.CanInt(v) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
	}

	strings := []any{
		"42",
		myString("42"),
		[]byte("42"),
		[]rune("42"),
		myBytes("42"),
		myRunes("42"),
		my42Marshaler{},
	}
	for _, v := range strings {
		v := v
		t.Run(fmt.Sprintf("WithConvertStrings + %T(%v)", v, v), func(t *testing.T) {
			if !into.CanInt(v, into.WithConvertStrings()) {
				t.Error("failed but should have succeeded for value:", v)
			}
		})
		t.Run(fmt.Sprintf("no convert + %T(%v)", v, v), func(t *testing.T) {
			if into.CanInt(v) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
	}

	badStrings := []any{
		"",
		new(string),
		new(myString),
		"foo",
		my42Marshaler{err: myError},
		myMarshaler{},
		myRunes(""),
		myBytes(""),
	}
	for _, v := range badStrings {
		v := v
		t.Run(fmt.Sprintf("strict string %T(%v)", v, v), func(t *testing.T) {
			if into.CanInt(v, into.WithConvertStrings()) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
		t.Run(fmt.Sprintf("WithoutMarshalerCheck + %T(%v)", v, v), func(t *testing.T) {
			if !into.CanInt(v, into.WithConvertStrings(), into.WithoutMarshalerCheck()) {
				t.Error("failed but should have succeeded for value:", v)
			}
		})
	}

	t.Run("WithConvertStrings disabled", func(t *testing.T) {
		if into.CanInt("123") {
			t.Error("unexpected success")
		}
	})
}

func TestIntInvalidFallback(t *testing.T) {
	t.Parallel()
	err := into.Try(func() {
		into.Int(nil, into.WithFallback("bad"))
	})
	if !strings.Contains(err.Error(), "invalid fallback") {
		t.Error("unexpected error (panic):", err)
	}
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := 42
		got := into.Int(42)
		if want != got {
			b.Fatal("bad result. want:", want, "got:", got)
		}
	}
}

func BenchmarkIntWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := 32
		got := into.Int(32, into.WithConvertStrings(), into.WithFallback(int(42)))
		if want != got {
			b.Fatal("bad result. want:", want, "got:", got)
		}
	}
}

func BenchmarkIntFallback(b *testing.B) {
	fallback := into.WithFallback(42)
	for i := 0; i < b.N; i++ {
		want := 42
		got := into.Int(nil, fallback)
		if want != got {
			b.Error("bad result. want:", want, "got:", got)
		}
	}
}
