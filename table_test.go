package into_test

import (
	"testing"

	"github.com/guregu/into"
)

type table[T comparable] []struct {
	name  string
	input any
	want  T
	opts  []into.Option
	err   error
}

func (tab table[T]) Run(t *testing.T, do func(any, ...into.Option) T) {
	t.Helper()
	for _, test := range tab {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			var got T
			err := into.Try(func() {
				got = do(test.input, test.opts...)
			})
			if err != nil && test.err == nil {
				t.Fatal("unexpected error (panic):", err)
			}
			if err == nil && test.err != nil {
				t.Error("expected error (panic), but did not see one")
			}
			if got != test.want {
				t.Error("bad return value. want:", test.want, "got:", got)
			}
		})
	}
}
