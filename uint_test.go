package into_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/guregu/into"
)

type myUint uint

func TestUint(t *testing.T) {
	t.Parallel()
	tests := table[uint]{
		{
			name:  "uint",
			input: uint(42),
			want:  42,
		},
		// {
		// 	name:  "uint without reflection",
		// 	input: myUint(42),
		// 	want:  0,
		// 	opts:  []into.Option{into.WithoutReflection()},
		// },
		{
			name:  "uint64",
			input: uint64(42),
			want:  42,
		},
		{
			name:  "uint32",
			input: uint32(42),
			want:  42,
		},
		{
			name:  "uint16",
			input: uint16(42),
			want:  42,
		},
		{
			name:  "uint8",
			input: uint8(42),
			want:  42,
		},
		{
			name:  "*uint",
			input: into.Ptr(uint(42)),
			want:  42,
		},
		{
			name:  "*uint64",
			input: into.Ptr(uint64(42)),
			want:  42,
		},
		{
			name:  "*uint32",
			input: into.Ptr(uint32(42)),
			want:  42,
		},
		{
			name:  "*uint16",
			input: into.Ptr(uint16(42)),
			want:  42,
		},
		{
			name:  "*uint8",
			input: into.Ptr(uint8(42)),
			want:  42,
		},
		{
			name:  "*uint(nil)",
			input: (*uint)(nil),
			want:  0,
		},
		{
			name:  "*uint64(nil)",
			input: (*uint64)(nil),
			want:  0,
		},
		{
			name:  "*uint32(nil)",
			input: (*uint32)(nil),
			want:  0,
		},
		{
			name:  "*uint16(nil)",
			input: (*uint16)(nil),
			want:  0,
		},
		{
			name:  "*uint8(nil)",
			input: (*uint8)(nil),
			want:  0,
		},
		{
			name:  "subtype",
			input: myUint(42),
			want:  42,
		},
		{
			name:  "subtype pointer",
			input: into.Ptr(myUint(42)),
			want:  42,
		},
		{
			name:  "nil subtype pointer",
			input: (*myUint)(nil),
			want:  0,
		},
		{
			name:  "subtype without reflection",
			input: myUint(42),
			want:  0,
			opts:  []into.Option{into.WithoutReflection()},
			err:   into.ErrInvalid{Value: 42, Type: "uint"},
		},
		{
			name:  "fallback",
			input: nil,
			want:  42,
			opts:  []into.Option{into.WithFallback(uint(42))},
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
			err:   into.ErrInvalid{Value: "abc", Type: "uint", Cause: &strconv.NumError{Func: "Atoi", Num: "abc", Err: strconv.ErrSyntax}},
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
			err:   into.ErrInvalid{Value: myString("42"), Type: "uint"},
		},
		{
			name:  "invalid type",
			input: struct{}{},
			want:  0,
			opts:  []into.Option{into.WithoutReflection()},
			err:   into.ErrInvalid{Value: struct{}{}, Type: "uint"},
		},
	}
	tests.Run(t, into.Uint)
}

func TestCanUint(t *testing.T) {
	t.Parallel()

	good := []any{
		uint(42),
		uint64(42),
		uint32(42),
		uint16(42),
		uint8(42),
		myUint(42),
		new(uint),
		new(myUint),
		nil,
	}
	for _, v := range good {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if !into.CanUint(v, into.WithoutMarshalerCheck()) {
				t.Errorf("failed but should have succeeded for value: %v (type %T)", v, v)
			}
		})
	}

	bad := []any{
		float64(1),
		float32(1),
		complex(1, 1),
		struct{}{},
		[]int{1},
	}
	for _, v := range bad {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if into.CanUint(v) {
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
			if !into.CanUint(v, into.WithConvertStrings()) {
				t.Error("failed but should have succeeded for value:", v)
			}
		})
		t.Run(fmt.Sprintf("no convert + %T(%v)", v, v), func(t *testing.T) {
			if into.CanUint(v) {
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
		t.Run(fmt.Sprintf("bad string %T(%v)", v, v), func(t *testing.T) {
			if into.CanUint(v, into.WithConvertStrings()) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
		t.Run(fmt.Sprintf("WithoutMarshalerCheck + %T(%v)", v, v), func(t *testing.T) {
			if !into.CanUint(v, into.WithConvertStrings(), into.WithoutMarshalerCheck()) {
				t.Error("failed but should have succeeded for value:", v)
			}
		})
	}

	t.Run("WithConvertStrings disabled", func(t *testing.T) {
		if into.CanUint("123") {
			t.Error("unexpected success")
		}
	})
}

func TestUintInvalidFallback(t *testing.T) {
	t.Parallel()
	err := into.Try(func() {
		into.Uint(nil, into.WithFallback("bad"))
	})
	if !strings.Contains(err.Error(), "invalid fallback") {
		t.Error("unexpected error (panic):", err)
	}
}

func BenchmarkUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := uint(42)
		got := into.Uint(uint(42))
		if want != got {
			b.Fatal("bad result. want:", want, "got:", got)
		}
	}
}

func BenchmarkUintWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := uint(32)
		got := into.Uint(uint(32), into.WithConvertStrings(), into.WithFallback(uint(42)))
		if want != got {
			b.Fatal("bad result. want:", want, "got:", got)
		}
	}
}

func BenchmarkUintFallback(b *testing.B) {
	fallback := into.WithFallback(uint(42))
	for i := 0; i < b.N; i++ {
		want := uint(42)
		got := into.Uint(nil, fallback)
		if want != got {
			b.Error("bad result. want:", want, "got:", got)
		}
	}
}
