package main

import (
	"testing"

	"github.com/Atnuhs/atcoder-cui/go-acl/testlib"
)

func TestKruskal(t *testing.T) {
	tests := map[string]struct {
		n          int
		edges      []*FullEdge
		wantWeight int
		wantEdges  int
		shouldFail bool
	}{
		"simple triangle": {
			n: 3,
			edges: []*FullEdge{
				{0, 1, 1},
				{1, 2, 2},
				{0, 2, 3},
			},
			wantWeight: 3,
			wantEdges:  2,
			shouldFail: false,
		},
		"disconnected graph": {
			n: 4,
			edges: []*FullEdge{
				{0, 1, 1},
				{2, 3, 2},
			},
			wantWeight: -1,
			wantEdges:  0,
			shouldFail: true,
		},
		"single node": {
			n:          1,
			edges:      []*FullEdge{},
			wantWeight: 0,
			wantEdges:  0,
			shouldFail: false,
		},
		"complete graph": {
			n: 4,
			edges: []*FullEdge{
				{0, 1, 1},
				{0, 2, 4},
				{0, 3, 3},
				{1, 2, 2},
				{1, 3, 5},
				{2, 3, 6},
			},
			wantWeight: 6, // 1 + 2 + 3
			wantEdges:  3,
			shouldFail: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			weight, edges := Kruskal(tc.n, tc.edges)

			testlib.AclAssert(t, tc.wantWeight, weight)

			if tc.shouldFail {
				testlib.AclAssert(t, ([]*FullEdge)(nil), edges)
			} else {
				testlib.AclAssert(t, tc.wantEdges, len(edges))

				// Verify it's actually a tree (n-1 edges for n nodes)
				if tc.n > 1 {
					testlib.AclAssert(t, tc.n-1, len(edges))
				}
			}
		})
	}
}
