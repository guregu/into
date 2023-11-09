package into_test

import (
	"fmt"
	"testing"

	"github.com/guregu/into"
)

func TestPtr(t *testing.T) {
	t.Parallel()
	want := 42
	got := into.Ptr(want)
	if got == nil {
		t.Fatal("unexpected nil pointer")
	}
	if want != *got {
		t.Error("unexpected value, want:", want, "got:", *got)
	}
}

func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()
		want := 42
		ptr := &want
		got := into.Value(ptr)
		if want != got {
			t.Error("unexpected value, want:", want, "got:", got)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		want := 0
		var ptr *int
		got := into.Value(ptr)
		if want != got {
			t.Error("unexpected value, want:", want, "got:", got)
		}
	})
}

func TestValueOr(t *testing.T) {
	t.Parallel()

	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()
		want := 42
		ptr := &want
		got := into.ValueOr(ptr, 1337)
		if want != got {
			t.Error("unexpected value, want:", want, "got:", got)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		want := 42
		var ptr *int
		got := into.ValueOr(ptr, 42)
		if want != got {
			t.Error("unexpected value, want:", want, "got:", got)
		}
	})
}

func ExamplePtr() {
	np := into.Ptr(42)
	fmt.Println(*np)
	// Output: 42
}

func ExampleValue() {
	var np *int
	value := into.Value(np)
	fmt.Println(value)
	// Output: 0
}
