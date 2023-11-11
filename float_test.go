package into_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/guregu/into"
)

type myFloat float64

func TestFloat(t *testing.T) {
	t.Parallel()
	tests := table[float64]{
		{
			name:  "float64",
			input: float64(42.5),
			want:  42.5,
		},
		{
			name:  "float32",
			input: float32(42.5),
			want:  42.5,
		},
		{
			name:  "*float64",
			input: into.Ptr(float64(42.5)),
			want:  42.5,
		},
		{
			name:  "*float32",
			input: into.Ptr(float32(42.5)),
			want:  42.5,
		},
		{
			name:  "*float64(nil)",
			input: (*float64)(nil),
			want:  0,
		},
		{
			name:  "*float32(nil)",
			input: (*float32)(nil),
			want:  0,
		},
		{
			name:  "subtype",
			input: myFloat(42.5),
			want:  42.5,
		},
		{
			name:  "subtype pointer",
			input: into.Ptr(myFloat(42.5)),
			want:  42.5,
		},
		{
			name:  "nil subtype pointer",
			input: (*myFloat)(nil),
			want:  0,
		},
		{
			name:  "subtype without reflection",
			input: myFloat(42.5),
			want:  0,
			opts:  []into.Option{into.WithoutReflection()},
			err:   into.ErrInvalid{Value: 42.5, Type: "int"},
		},
		{
			name:  "fallback",
			input: nil,
			want:  42.5,
			opts:  []into.Option{into.WithFallback(42.5)},
		},
		{
			name:  "default fallback",
			input: nil,
			want:  0,
		},
		{
			name:  "string conversion",
			input: "42.5",
			want:  42.5,
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
			input: myString("42.5"),
			want:  42.5,
			opts:  []into.Option{into.WithConvertStrings()},
		},
		{
			name:  "string conversion using subtype without reflection",
			input: myString("42.5"),
			want:  0,
			opts:  []into.Option{into.WithConvertStrings(), into.WithoutReflection()},
			err:   into.ErrInvalid{Value: myString("42.5"), Type: "int"},
		},
		{
			name:  "invalid type",
			input: struct{}{},
			want:  0,
			opts:  []into.Option{into.WithoutReflection()},
			err:   into.ErrInvalid{Value: struct{}{}, Type: "int"},
		},
	}
	tests.Run(t, into.Float)
}

func TestCanFloat(t *testing.T) {
	t.Parallel()

	good := []any{
		float64(42.5),
		float32(42.5),
		myFloat(42.5),
		new(float64),
		new(float32),
		new(myFloat),
		nil,
	}
	for _, v := range good {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if !into.CanFloat(v) {
				t.Error("failed but should have succeeded for value:", v)
			}
		})
	}

	bad := []any{
		int(1),
		uint(1),
		complex(1, 1),
		struct{}{},
		[]int{1},
	}
	for _, v := range bad {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if into.CanFloat(v) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
	}

	strings := []any{
		"42.5",
		myString("42.5"),
		[]byte("42.5"),
		[]rune("42.5"),
		myBytes("42.5"),
		myRunes("42.5"),
		my42Marshaler{},
	}
	for _, v := range strings {
		v := v
		t.Run(fmt.Sprintf("WithConvertStrings + %T(%v)", v, v), func(t *testing.T) {
			if !into.CanFloat(v, into.WithConvertStrings()) {
				t.Error("failed but should have succeeded for value:", v)
			}
		})
		t.Run(fmt.Sprintf("no convert + %T(%v)", v, v), func(t *testing.T) {
			if into.CanFloat(v) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
	}

	badStrings := []any{
		[]int{1},
		"foo",
		"",
		new(string),
		new(myString),
		my42Marshaler{err: myError},
		myMarshaler{},
		myRunes(""),
		myBytes(""),
	}
	for _, v := range badStrings {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if into.CanFloat(v, into.WithConvertStrings()) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
	}

	t.Run("WithConvertStrings disabled", func(t *testing.T) {
		if into.CanFloat("123") {
			t.Error("unexpected success")
		}
	})
}

func TestFloatInvalidFallback(t *testing.T) {
	t.Parallel()
	err := into.Try(func() {
		into.Float(nil, into.WithFallback("bad"))
	})
	if !strings.Contains(err.Error(), "invalid fallback") {
		t.Error("unexpected error (panic):", err)
	}
}

func BenchmarkFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := 42.5
		got := into.Float(42.5)
		if want != got {
			b.Fatal("bad result. want:", want, "got:", got)
		}
	}
}

func BenchmarkFloatWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := float64(42.5)
		got := into.Float(42.5, into.WithConvertStrings(), into.WithFallback(float64(42.5)))
		if want != got {
			b.Fatal("bad result. want:", want, "got:", got)
		}
	}
}

func BenchmarkFloatFallback(b *testing.B) {
	fallback := into.WithFallback(42.5)
	for i := 0; i < b.N; i++ {
		want := 42.5
		got := into.Float(nil, fallback)
		if want != got {
			b.Error("bad result. want:", want, "got:", got)
		}
	}
}
