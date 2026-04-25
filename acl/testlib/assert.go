package testlib

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func AclAssert(t *testing.T, want, got any) {
	t.Helper()
	AclAssertWithOpts(t, want, got)
}

func AclAssertWithOpts(t *testing.T, want any, got any, opts ...cmp.Option) {
	t.Helper()
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Errorf("want(-), got(+)\n%s\n", diff)
	}
}

func AclAssertEquateEmpty(t *testing.T, want, got any) {
	t.Helper()
	AclAssertWithOpts(t, want, got, cmpopts.EquateEmpty())
}
