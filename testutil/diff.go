package testutil

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertObject(t *testing.T, got any, want any) error {
	t.Helper()
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		return fmt.Errorf("result is difference: %v", diff)
	}
	return nil
}
