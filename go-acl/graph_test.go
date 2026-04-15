package main

import (
	"testing"

	"github.com/Atnuhs/atcoder-cui/go-acl/testlib"
)

func TestNewWEdge(t *testing.T) {
	tests := map[string]struct {
		from, to, weight int
		want             *FullEdge
	}{
		"simple edge": {
			from:   0,
			to:     1,
			weight: 5,
			want:   &FullEdge{From: 0, To: 1, Weight: 5},
		},
		"negative weight": {
			from:   2,
			to:     3,
			weight: -1,
			want:   &FullEdge{From: 2, To: 3, Weight: -1},
		},
		"zero weight": {
			from:   1,
			to:     2,
			weight: 0,
			want:   &FullEdge{From: 1, To: 2, Weight: 0},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := &FullEdge{tc.from, tc.to, tc.weight}
			testlib.AclAssert(t, tc.want, got)
		})
	}
}

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
				&FullEdge{0, 1, 1},
				&FullEdge{1, 2, 2},
				&FullEdge{0, 2, 3},
			},
			wantWeight: 3,
			wantEdges:  2,
			shouldFail: false,
		},
		"disconnected graph": {
			n: 4,
			edges: []*FullEdge{
				&FullEdge{0, 1, 1},
				&FullEdge{2, 3, 2},
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
				&FullEdge{0, 1, 1},
				&FullEdge{0, 2, 4},
				&FullEdge{0, 3, 3},
				&FullEdge{1, 2, 2},
				&FullEdge{1, 3, 5},
				&FullEdge{2, 3, 6},
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
