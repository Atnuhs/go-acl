package main

import (
	"testing"

	"github.com/Atnuhs/atcoder-cui/go-acl/testlib"
)

func TestNewIntAbel(t *testing.T) {
	ab := NewIntAbel()
	testlib.AclAssert(t, 0, ab.E)
	testlib.AclAssert(t, 5, ab.Add(2, 3))
	testlib.AclAssert(t, -1, ab.Sub(2, 3))
	testlib.AclAssert(t, -5, ab.Minus(5))
}

func TestNewFloat64Abel(t *testing.T) {
	ab := NewFloat64Abel()
	testlib.AclAssert(t, 0.0, ab.E)
	testlib.AclAssert(t, 5.5, ab.Add(2.5, 3.0))
	testlib.AclAssert(t, -0.5, ab.Sub(2.5, 3.0))
	testlib.AclAssert(t, -2.5, ab.Minus(2.5))
}

func TestNewXorAbel(t *testing.T) {
	ab := NewXorAbel()
	testlib.AclAssert(t, 0, ab.E)
	testlib.AclAssert(t, 0b0110, ab.Add(0b0101, 0b0011)) // 5 ^ 3 = 6
	testlib.AclAssert(t, 0b0110, ab.Sub(0b0101, 0b0011)) // XOR: Sub = Add
	testlib.AclAssert(t, 5, ab.Minus(5))                 // XOR: 逆元は自分自身
}

func TestNewModAbel(t *testing.T) {
	ab := NewModAbel(7)
	testlib.AclAssert(t, 0, ab.E)
	testlib.AclAssert(t, 1, ab.Add(5, 3))  // (5+3) % 7 = 1
	testlib.AclAssert(t, 2, ab.Sub(5, 3))  // (5-3) % 7 = 2
	testlib.AclAssert(t, 5, ab.Sub(3, 5))  // (3-5) % 7 = -2 % 7 = 5
	testlib.AclAssert(t, 4, ab.Minus(3))   // (0-3) % 7 = 4
}

func TestNewPartialUF(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(5, ab)

	// Initially, all elements should have size 1 and weight 0
	for i := 0; i < 5; i++ {
		testlib.AclAssert(t, 1, uf.Size(i))
		testlib.AclAssert(t, i, uf.Root(i))
		testlib.AclAssert(t, 0, uf.Weight(i))
	}
}

func TestPartialUF_BasicUnion(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(5, ab)

	// Union(x, y, w) means: weight[y] - weight[x] = w
	// So weight[y] = weight[x] + w
	uf.Union(0, 1, 10) // weight[1] - weight[0] = 10

	testlib.AclAssert(t, true, uf.Family(0, 1))
	testlib.AclAssert(t, 2, uf.Size(0))
	testlib.AclAssert(t, 2, uf.Size(1))

	// Diff(x, y) = Weight(x) - Weight(y) = weight[x] - weight[y]
	// So Diff(1, 0) = weight[1] - weight[0] = 10
	testlib.AclAssert(t, 10, uf.Diff(1, 0))
	testlib.AclAssert(t, -10, uf.Diff(0, 1))
}

func TestPartialUF_ChainUnion(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(5, ab)

	// weight[1] - weight[0] = 5
	// weight[2] - weight[1] = 3
	// => weight[2] - weight[0] = 8
	uf.Union(0, 1, 5)
	uf.Union(1, 2, 3)

	testlib.AclAssert(t, true, uf.Family(0, 2))
	testlib.AclAssert(t, 5, uf.Diff(1, 0))
	testlib.AclAssert(t, 3, uf.Diff(2, 1))
	testlib.AclAssert(t, 8, uf.Diff(2, 0))
}

func TestPartialUF_MergeGroups(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(6, ab)

	// Create two separate groups
	// Group 1: 0 - 1 - 2
	// weight[1] - weight[0] = 10
	// weight[2] - weight[1] = 20
	uf.Union(0, 1, 10)
	uf.Union(1, 2, 20)

	// Group 2: 3 - 4 - 5
	// weight[4] - weight[3] = 5
	// weight[5] - weight[4] = 7
	uf.Union(3, 4, 5)
	uf.Union(4, 5, 7)

	// Verify groups are separate
	testlib.AclAssert(t, false, uf.Family(0, 3))

	// Merge groups: weight[3] - weight[2] = 100
	uf.Union(2, 3, 100)

	testlib.AclAssert(t, true, uf.Family(0, 5))
	testlib.AclAssert(t, 6, uf.Size(0))

	// Verify all diffs
	// weight[2] - weight[0] = 30
	testlib.AclAssert(t, 30, uf.Diff(2, 0))
	// weight[3] - weight[2] = 100
	testlib.AclAssert(t, 100, uf.Diff(3, 2))
	// weight[5] - weight[3] = 12
	testlib.AclAssert(t, 12, uf.Diff(5, 3))
	// weight[5] - weight[0] = 30 + 100 + 12 = 142
	testlib.AclAssert(t, 142, uf.Diff(5, 0))
}

func TestPartialUF_UnionSameElement(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(3, ab)

	uf.Union(0, 0, 10)
	testlib.AclAssert(t, 1, uf.Size(0))
	testlib.AclAssert(t, 0, uf.Weight(0))
}

func TestPartialUF_UnionAlreadyConnected(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(3, ab)

	uf.Union(0, 1, 10)
	uf.Union(0, 1, 20) // Should not change anything

	testlib.AclAssert(t, 10, uf.Diff(1, 0))
}

func TestPartialUF_NegativeWeights(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(3, ab)

	// weight[1] - weight[0] = -5
	uf.Union(0, 1, -5)

	testlib.AclAssert(t, -5, uf.Diff(1, 0))
	testlib.AclAssert(t, 5, uf.Diff(0, 1))
}

func TestPartialUF_PathCompression(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(10, ab)

	// Create a long chain: 0 -> 1 -> 2 -> ... -> 9
	for i := 0; i < 9; i++ {
		uf.Union(i, i+1, i+1) // weight[i+1] - weight[i] = i+1
	}

	// weight[9] - weight[0] = 1+2+3+...+9 = 45
	testlib.AclAssert(t, 45, uf.Diff(9, 0))

	// After path compression, queries should still work
	testlib.AclAssert(t, 45, uf.Diff(9, 0))
	testlib.AclAssert(t, 15, uf.Diff(5, 0)) // 1+2+3+4+5 = 15
}

func TestPartialUF_UnionBySize(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(10, ab)

	// Create a large group (0-4)
	for i := 0; i < 4; i++ {
		uf.Union(i, i+1, 1)
	}
	testlib.AclAssert(t, 5, uf.Size(0))

	// Create a small group (5-6)
	uf.Union(5, 6, 2)
	testlib.AclAssert(t, 2, uf.Size(5))

	// Merge: weight[6] - weight[4] = 10
	// Small group should be merged into large group
	uf.Union(4, 6, 10)

	testlib.AclAssert(t, 7, uf.Size(0))
	testlib.AclAssert(t, true, uf.Family(0, 6))

	// Verify weights are correct after merge
	// weight[4] - weight[0] = 4
	// weight[6] - weight[4] = 10
	// weight[6] - weight[5] = 2
	testlib.AclAssert(t, 4, uf.Diff(4, 0))
	testlib.AclAssert(t, 10, uf.Diff(6, 4))
	testlib.AclAssert(t, 2, uf.Diff(6, 5))
}

func TestPartialUF_ReverseUnion(t *testing.T) {
	ab := NewIntAbel()
	uf := NewPartialUF(5, ab)

	// Create groups with different sizes to test swap in Union
	// Large group: 0, 1, 2
	uf.Union(0, 1, 10)
	uf.Union(1, 2, 20)

	// Small group: 3
	// Merge: weight[0] - weight[3] = 5
	// This triggers the swap branch because Size(r3) < Size(r0)
	uf.Union(3, 0, 5)

	testlib.AclAssert(t, true, uf.Family(0, 3))
	testlib.AclAssert(t, 5, uf.Diff(0, 3))
	testlib.AclAssert(t, 15, uf.Diff(1, 3)) // 5 + 10
	testlib.AclAssert(t, 35, uf.Diff(2, 3)) // 5 + 10 + 20
}
