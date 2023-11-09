package into_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/guregu/into"
)

func TestTry(t *testing.T) {
	t.Parallel()
	t.Run("panic(error)", func(t *testing.T) {
		t.Parallel()
		want := fmt.Errorf("test error")
		got := into.Try(func() {
			panic(want)
		})
		if want != got {
			t.Error("bad error. want:", want, "got:", got)
		}
	})

	t.Run("panic(non-error)", func(t *testing.T) {
		t.Parallel()
		want := into.Panic{Value: "abcdef"}
		got := into.Try(func() {
			panic(want.Value)
		})
		if want != got {
			t.Error("bad error. want:", want, "got:", got)
		}
		if !strings.Contains(got.Error(), "panic:") {
			t.Error("unexpected error message:", got.Error())
		}
	})

	t.Run("no panic", func(t *testing.T) {
		got := into.Try(func() {})
		if got != nil {
			t.Error("unexpected error:", got)
		}
	})
}

func ExampleMaybe() {
	_, err := into.Maybe(into.Int, "cat")
	fmt.Println(err)
	// Output: into: value cat of type string is not a int
}

func TestMaybe(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		_, err := into.Maybe(into.Int, 123)
		if err != nil {
			t.Error("unexpected error")
		}
	})
	t.Run("with error", func(t *testing.T) {
		_, err := into.Maybe(into.Int, "blah")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func BenchmarkMaybeSuccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := into.Maybe(into.Int, 123)
		if err != nil {
			b.Fatal("unexpected error")
		}
	}
}

func BenchmarkMaybeWithError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := into.Maybe(into.Int, "blah")
		if err == nil {
			b.Fatal("expected error")
		}
	}
}

func BenchmarkStrconvError(b *testing.B) {
	var x any = "abcdsjfb"
	for i := 0; i < b.N; i++ {
		_, err := strconv.Atoi(x.(string))
		if err == nil {
			b.Fatal("expected error")
		}
	}
}
