package into_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/guregu/into"
)

func TestErrInvalid(t *testing.T) {
	err := into.ErrInvalid{Value: "foo", Type: "int", Cause: myError}
	if !errors.Is(err, myError) {
		t.Error("error wrapping failed")
	}
	if !strings.Contains(err.Error(), myError.Error()) {
		t.Error("bad error message")
	}
}
