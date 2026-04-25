package acl

import (
	"testing"

	"github.com/Atnuhs/go-acl/acl/testlib"
)

func TestNewUnionFind(t *testing.T) {
	tests := map[string]struct {
		n    int
		want int
	}{
		"size 1": {n: 1, want: 1},
		"size 5": {n: 5, want: 5},
		"size 10": {n: 10, want: 10},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			uf := NewUnionFind(tc.n)
			testlib.AclAssert(t, tc.want, len(uf.data))
			
			// Initially, all elements should have size 1
			for i := 0; i < tc.n; i++ {
				testlib.AclAssert(t, 1, uf.Size(i))
				testlib.AclAssert(t, i, uf.Root(i))
			}
		})
	}
}

func TestUnionFind_Root(t *testing.T) {
	uf := NewUnionFind(5)
	
	// Initially, each element is its own root
	for i := 0; i < 5; i++ {
		testlib.AclAssert(t, i, uf.Root(i))
	}
	
	// After union, root should change
	uf.Union(0, 1)
	root0 := uf.Root(0)
	root1 := uf.Root(1)
	testlib.AclAssert(t, root0, root1)
}

func TestUnionFind_Family(t *testing.T) {
	uf := NewUnionFind(5)
	
	// Initially, no elements are in the same family
	testlib.AclAssert(t, false, uf.Family(0, 1))
	testlib.AclAssert(t, false, uf.Family(0, 2))
	
	// After union, elements should be in the same family
	uf.Union(0, 1)
	testlib.AclAssert(t, true, uf.Family(0, 1))
	testlib.AclAssert(t, false, uf.Family(0, 2))
	
	uf.Union(1, 2)
	testlib.AclAssert(t, true, uf.Family(0, 1))
	testlib.AclAssert(t, true, uf.Family(0, 2))
	testlib.AclAssert(t, true, uf.Family(1, 2))
}

func TestUnionFind_Size(t *testing.T) {
	uf := NewUnionFind(5)
	
	// Initially, all components have size 1
	for i := 0; i < 5; i++ {
		testlib.AclAssert(t, 1, uf.Size(i))
	}
	
	// After union, size should increase
	uf.Union(0, 1)
	testlib.AclAssert(t, 2, uf.Size(0))
	testlib.AclAssert(t, 2, uf.Size(1))
	testlib.AclAssert(t, 1, uf.Size(2))
	
	uf.Union(0, 2)
	testlib.AclAssert(t, 3, uf.Size(0))
	testlib.AclAssert(t, 3, uf.Size(1))
	testlib.AclAssert(t, 3, uf.Size(2))
}

func TestUnionFind_Union(t *testing.T) {
	uf := NewUnionFind(5)
	
	// Union same element should not change anything
	uf.Union(0, 0)
	testlib.AclAssert(t, 1, uf.Size(0))
	
	// Union different elements
	uf.Union(0, 1)
	testlib.AclAssert(t, true, uf.Family(0, 1))
	testlib.AclAssert(t, 2, uf.Size(0))
	
	// Union already connected elements
	uf.Union(0, 1)
	testlib.AclAssert(t, 2, uf.Size(0))
	
	// Chain unions
	uf.Union(2, 3)
	uf.Union(0, 2)
	testlib.AclAssert(t, true, uf.Family(0, 3))
	testlib.AclAssert(t, 4, uf.Size(0))
}