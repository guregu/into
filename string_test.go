package into_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/guregu/into"
)

type myString string
type myBytes []byte
type myRunes []rune

func TestString(t *testing.T) {
	t.Parallel()
	tests := table[string]{
		{
			name:  "string",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "*string",
			input: into.Ptr("hello"),
			want:  "hello",
		},
		{
			name:  "*string(nil)",
			input: (*string)(nil),
			want:  "",
		},
		{
			name:  "subtype",
			input: myString("hello"),
			want:  "hello",
		},
		{
			name:  "subtype without reflection",
			input: myString("hello"),
			want:  "",
			opts:  []into.Option{into.WithoutReflection()},
			err:   into.ErrInvalid{Value: 42, Type: "string"},
		},
		{
			name:  "[]byte",
			input: []byte("hello"),
			want:  "hello",
		},
		{
			name:  "[]byte(nil)",
			input: []byte(nil),
			want:  "",
		},
		{
			name:  "bytes subtype",
			input: myBytes("hello"),
			want:  "hello",
		},
		{
			name:  "bytes subtype nil",
			input: myBytes(nil),
			want:  "",
		},
		{
			name:  "bytes subtype pointer",
			input: into.Ptr(myBytes("hello")),
			want:  "hello",
		},
		{
			name:  "bytes subtype pointer but nil",
			input: (*myBytes)(nil),
			want:  "",
		},
		{
			name:  "[]int",
			input: []int{1},
			want:  "",
			err:   into.ErrInvalid{Value: []int{1}, Type: "string"},
		},
		{
			name:  "rune",
			input: rune('A'),
			want:  "A",
		},
		{
			name:  "*rune",
			input: into.Ptr(rune('A')),
			want:  "A",
		},
		{
			name:  "*rune(nil)",
			input: (*rune)(nil),
			want:  "",
		},
		{
			name:  "[]rune",
			input: []rune("hello"),
			want:  "hello",
		},
		{
			name:  "[]rune(nil)",
			input: []rune(nil),
			want:  "",
		},
		{
			name:  "runes subtype",
			input: myRunes("hello"),
			want:  "hello",
		},
		{
			name:  "runes subtype nil",
			input: myRunes(nil),
			want:  "",
		},
		{
			name:  "runes subtype pointer",
			input: into.Ptr(myRunes("hello")),
			want:  "hello",
		},
		{
			name:  "runes subtype pointer but nil",
			input: (*myRunes)(nil),
			want:  "",
		},
		{
			name:  "TextMarshaler",
			input: myMarshaler{},
			want:  "hello",
		},
		{
			name:  "TextMarshaler with error",
			input: myMarshaler{err: myError},
			want:  "",
			err:   into.ErrInvalid{Value: struct{}{}, Type: "string", Cause: myError},
		},
		{
			name:  "Stringer",
			input: myStringer{},
			want:  "hello",
		},
		// {
		// 	name:  "*[]rune",
		// 	input: into.Ptr([]rune("hello")),
		// 	want:  "hello",
		// },
		// {
		// 	name:  "*[]rune(nil)",
		// 	input: (*[]rune)(nil),
		// 	want:  "",
		// },
		{
			name:  "fallback",
			input: nil,
			want:  "hello",
			opts:  []into.Option{into.WithFallback("hello")},
		},
		{
			name:  "default fallback",
			input: nil,
			want:  "",
		},
		{
			name:  "invalid type",
			input: struct{}{},
			want:  "",
			opts:  []into.Option{into.WithoutReflection()},
			err:   into.ErrInvalid{Value: struct{}{}, Type: "string"},
		},
	}
	tests.Run(t, into.String)
}

func TestStringInvalidFallback(t *testing.T) {
	t.Parallel()
	err := into.Try(func() {
		into.String(nil, into.WithFallback(1337))
	})
	if !strings.Contains(err.Error(), "invalid fallback") {
		t.Error("unexpected error (panic):", err)
	}
}

func TestCanString(t *testing.T) {
	t.Parallel()

	good := []any{
		"abc",
		myString("abc"),
		[]rune("abc"),
		rune('a'),
		new(string),
		new(rune),
		new(myString),
		myBytes("abc"),
		myMarshaler{},
		myStringer{},
		nil,
	}

	bad := []any{
		float64(1),
		float32(1),
		complex(1, 1),
		struct{}{},
	}

	strictlyBad := []any{
		myMarshaler{err: myError},
	}

	for _, v := range good {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if !into.CanString(v) {
				t.Error("failed but should have succeeded for value:", v)
			}
		})
	}

	for _, v := range bad {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if into.CanString(v) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
	}

	for _, v := range strictlyBad {
		v := v
		t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
			if into.CanString(v) {
				t.Error("succeeded but should have failed for value:", v)
			}
		})
	}
}

var myError = fmt.Errorf("test error")

type myMarshaler struct {
	err error
}

func (m myMarshaler) MarshalText() ([]byte, error) {
	return []byte("hello"), m.err
}

type my42Marshaler struct {
	err error
}

func (m my42Marshaler) MarshalText() ([]byte, error) {
	return []byte("42"), m.err
}

type myStringer struct{}

func (myStringer) String() string {
	return "hello"
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := "hello"
		got := into.String("hello")
		if want != got {
			b.Fatal("bad result. want:", want, "got:", got)
		}
	}
}

func BenchmarkStringWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := "hello"
		got := into.String("hello", into.WithConvertStrings(), into.WithFallback("abc"))
		if want != got {
			b.Fatal("bad result. want:", want, "got:", got)
		}
	}
}

func BenchmarkStringFallback(b *testing.B) {
	fallback := into.WithFallback("hello")
	for i := 0; i < b.N; i++ {
		want := "hello"
		got := into.String(nil, fallback)
		if want != got {
			b.Error("bad result. want:", want, "got:", got)
		}
	}
}
