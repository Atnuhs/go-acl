package acl

import (
	"fmt"
	"testing"

	"github.com/Atnuhs/go-acl/acl/testlib"
)

func TestNewPair(t *testing.T) {
	tests := map[string]struct {
		u, v any
		want *Pair[any, any]
	}{
		"int pair": {
			u:    1,
			v:    2,
			want: &Pair[any, any]{1, 2},
		},
		"string pair": {
			u:    "hello",
			v:    "world",
			want: &Pair[any, any]{"hello", "world"},
		},
		"mixed pair": {
			u:    42,
			v:    "test",
			want: &Pair[any, any]{42, "test"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewPair(tc.u, tc.v)
			testlib.AclAssert(t, tc.want.U, got.U)
			testlib.AclAssert(t, tc.want.V, got.V)
		})
	}
}

func TestPair_String(t *testing.T) {
	tests := map[string]struct {
		setupPair func() fmt.Stringer
		want      string
	}{
		"int pair": {
			setupPair: func() fmt.Stringer { return NewPair(1, 2) },
			want:      "1 2",
		},
		"string pair": {
			setupPair: func() fmt.Stringer { return NewPair("hello", "world") },
			want:      "hello world",
		},
		"mixed pair": {
			setupPair: func() fmt.Stringer { return NewPair(42, "test") },
			want:      "42 test",
		},
		"negative numbers": {
			setupPair: func() fmt.Stringer { return NewPair(-1, -2) },
			want:      "-1 -2",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pair := tc.setupPair()
			got := pair.String()
			testlib.AclAssert(t, tc.want, got)
		})
	}
}
