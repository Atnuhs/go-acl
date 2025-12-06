package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Atnuhs/atcoder-cui/go-acl/testlib"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/constraints"
)

func Test_lIdx(t *testing.T) {
	tests := map[string]struct {
		idx  int
		want int
	}{
		"odd":  {idx: 5, want: 4},
		"even": {idx: 6, want: 6},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := lIdx(tc.idx)
			if tc.want != got {
				t.Errorf("want %d but got %d", tc.want, got)
			}
		})
	}
}

func Fuzz_lIdx(f *testing.F) {
	f.Add(0)
	f.Add(INF)

	f.Fuzz(func(t *testing.T, idx int) {
		if idx < 0 {
			return
		}

		want := (idx / 2) * 2
		got := lIdx(idx)
		if want != got {
			t.Errorf("idx: %d, want: %d, got: %d", idx, want, got)
		}
	})
}

func Test_rIdx(t *testing.T) {
	tests := map[string]struct {
		idx  int
		want int
	}{
		"odd":  {idx: 5, want: 5},
		"even": {idx: 6, want: 7},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := rIdx(tc.idx)
			if tc.want != got {
				t.Errorf("want %d but got %d", tc.want, got)
			}
		})
	}
}

func Fuzz_rIdx(f *testing.F) {
	f.Add(0)
	f.Add(INF)

	f.Fuzz(func(t *testing.T, idx int) {
		if idx < 0 {
			return
		}

		want := (idx/2)*2 + 1
		got := rIdx(idx)
		if want != got {
			t.Errorf("idx: %d, want: %d, got: %d", idx, want, got)
		}
	})
}

func Test_pIdx(t *testing.T) {
	tests := []struct {
		idx  int
		want int
	}{
		{2, 0},
		{3, 0},
		{4, 0},
		{5, 0},
		{6, 2},
		{7, 2},
		{8, 2},
		{9, 2},
		{10, 4},
		{13, 4},
		{14, 6},
		{17, 6},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			got := pIdx(tc.idx)
			if tc.want != got {
				t.Errorf("want %d but got %d", tc.want, got)
			}
		})
	}
}

func Fuzz_pIdx(f *testing.F) {
	f.Add(2)
	f.Add(3)
	f.Add(INF)

	f.Fuzz(func(t *testing.T, idx int) {
		if idx < 2 || idx > INF {
			return
		}

		want := ((idx - 2) / 4) * 2
		got := pIdx(idx)
		if want != got {
			t.Errorf("idx: %d, want: %d, got: %d", idx, want, got)
		}
	})
}

func Test_cIdx(t *testing.T) {
	tests := map[string]struct {
		idx  int
		want int
	}{
		"root0": {idx: 0, want: 2},
		"root1": {idx: 1, want: 3},
		"c1":    {idx: 2, want: 6},
		"c2":    {idx: 3, want: 7},
		"c3":    {idx: 4, want: 10},
		"c4":    {idx: 5, want: 11},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := cIdx(tc.idx)
			if tc.want != got {
				t.Errorf("want %d but got %d", tc.want, got)
			}
		})
	}
}

func Fuzz_cIdx(f *testing.F) {
	f.Add(2)
	f.Add(3)
	f.Add(INF)

	f.Fuzz(func(t *testing.T, idx int) {
		if idx < 2 || idx > INF {
			return
		}

		want := 2*(idx+1) - idx%2
		got := cIdx(idx)
		if want != got {
			t.Errorf("idx: %d, want: %d, got: %d", idx, want, got)
		}
	})
}

func check_heap[T constraints.Ordered](t *testing.T, pq *DEPQ[T]) {
	t.Helper()
	for i := range pq.values {
		cl := cIdx(i)
		cr := cl + 2

		if cl >= pq.Size() {
			continue
		}

		if i&1 == 1 {
			// min heap
			if pq.values[cl] < pq.values[i] {
				t.Errorf("parent(ID,value)=(%d,%v), child(ID,value)=(%d, %v) at min heap", i, pq.values[i], cl, pq.values[cl])
			}
		} else {
			// max heap
			if pq.values[cl] > pq.values[i] {
				t.Errorf("parent(ID,value)=(%d,%v), child(ID,value)=(%d, %v) at max heap", i, pq.values[i], cl, pq.values[cl])
			}
		}

		if cr >= pq.Size() {
			continue
		}

		if i&1 == 1 {
			// min heap
			if pq.values[cr] < pq.values[i] {
				t.Errorf("parent(ID,value)=(%d,%v), child(ID,value)=(%d, %v) at min heap", i, pq.values[i], cr, pq.values[cr])
			}
		} else {
			// max heap
			if pq.values[cr] > pq.values[i] {
				t.Errorf("parent(ID,value)=(%d,%v), child(ID,value)=(%d, %v) at max heap", i, pq.values[i], cr, pq.values[cr])
			}
		}
	}
}

func TestDEPQ_Push(t *testing.T) {
	tests := map[string]struct {
		values []int
	}{
		"0~9": {values: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pq := NewDEPQ[int]()
			for _, v := range tc.values {
				pq.Push(v)
			}
			t.Log(pq)
			check_heap(t, pq)
		})
	}
}

func TestDEPQ_Constructor(t *testing.T) {
	tests := map[string]struct {
		values []int
	}{
		"0~9": {values: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pq := NewDEPQ[int](tc.values...)
			t.Log(pq)
			check_heap(t, pq)
		})
	}
}

func TestDEPQ_PopMax(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := make([]int, len(values))
	copy(want, values)
	for i, j := 0, len(want)-1; i < j; i, j = i+1, j-1 {
		want[i], want[j] = want[j], want[i]
	}
	pq := NewDEPQ(values...)
	got := make([]int, 0, len(want))
	for !pq.Empty() {
		got = append(got, pq.PopMax())
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want(-), got(+)\n%s", diff)
	}
}

func TestDEPQ_PopMin(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := make([]int, len(values))
	copy(want, values)
	pq := NewDEPQ(values...)
	got := make([]int, 0, len(want))
	for !pq.Empty() {
		got = append(got, pq.PopMin())
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want(-), got(+)\n%s", diff)
	}
}

func TestDEPQ_EmptyOperations(t *testing.T) {
	tests := map[string]struct {
		setup func() *DEPQ[int]
		want  bool
	}{
		"empty queue": {
			setup: func() *DEPQ[int] { return NewDEPQ[int]() },
			want:  true,
		},
		"non-empty queue": {
			setup: func() *DEPQ[int] { return NewDEPQ(1, 2, 3) },
			want:  false,
		},
		"queue after popping all": {
			setup: func() *DEPQ[int] {
				pq := NewDEPQ(1)
				pq.PopMin()
				return pq
			},
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pq := tc.setup()
			got := pq.Empty()
			testlib.AclAssert(t, tc.want, got)
		})
	}
}

func TestDEPQ_SizeOperations(t *testing.T) {
	tests := map[string]struct {
		setup func() *DEPQ[int]
		want  int
	}{
		"empty": {
			setup: func() *DEPQ[int] { return NewDEPQ[int]() },
			want:  0,
		},
		"single element": {
			setup: func() *DEPQ[int] { return NewDEPQ(42) },
			want:  1,
		},
		"multiple elements": {
			setup: func() *DEPQ[int] { return NewDEPQ(1, 2, 3, 4, 5) },
			want:  5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pq := tc.setup()
			got := pq.Size()
			testlib.AclAssert(t, tc.want, got)
		})
	}
}

func TestDEPQ_GetOperations(t *testing.T) {
	pq := NewDEPQ(5, 1, 9, 3, 7)

	testlib.AclAssert(t, 1, pq.GetMin())
	testlib.AclAssert(t, 9, pq.GetMax())

	// Verify getting doesn't modify the queue
	testlib.AclAssert(t, 5, pq.Size())
	testlib.AclAssert(t, 1, pq.GetMin())
	testlib.AclAssert(t, 9, pq.GetMax())
}

func FuzzNewDEPQ(f *testing.F) {
	f.Add(10)
	f.Add(1_000_000)
	f.Fuzz(func(t *testing.T, a int) {
		if a <= 0 || 1_000_000 < a {
			return
		}
		arr := L1[int](a)
		for i := range arr {
			arr[i] = rand.Intn(INF)
		}
		dpq := NewDEPQ(arr...)
		check_heap(t, dpq)
	})
}

func FuzzDEPQ(f *testing.F) {
	f.Add(10)
	f.Add(1_000_000)
	f.Fuzz(func(t *testing.T, n int) {
		if n <= 0 || 1_000 < n {
			return
		}
		arr := L1[int](n)
		for i := range arr {
			arr[i] = rand.Intn(1_000)
		}
		mi, ma := arr[0], arr[0]
		dpq := NewDEPQ[int]()

		// Insert to heap test
		for _, v := range arr {
			dpq.Push(v)
			check_heap(t, dpq)
			mi = Min(mi, v)
			ma = Max(ma, v)

			if gmi := dpq.GetMin(); gmi != mi {
				t.Fatalf("want %d but got %d", mi, gmi)
			}
			if gma := dpq.GetMax(); gma != ma {
				t.Fatalf("want %d but got %d", ma, gma)
			}
		}

		// PopTest
		for !dpq.Empty() {
			q := rand.Intn(2)
			switch q {
			case 0:
				gma := dpq.PopMax()
				if gma > ma {
					t.Fatalf("shoudl be <= ma: %d but got %d", ma, gma)
				}
				ma = gma
			case 1:
				gmi := dpq.PopMin()
				if gmi < mi {
					t.Fatalf("should be >= mi: %d but got %d", mi, gmi)
				}
				mi = gmi
			}
			check_heap(t, dpq)
		}
	})
}

func TestDEPQ_EdgeCasesOperations(t *testing.T) {
	t.Run("single element operations", func(t *testing.T) {
		pq := NewDEPQ(42)
		testlib.AclAssert(t, 42, pq.GetMin())
		testlib.AclAssert(t, 42, pq.GetMax())

		min := pq.PopMin()
		testlib.AclAssert(t, 42, min)
		testlib.AclAssert(t, true, pq.Empty())
	})

	t.Run("two element operations", func(t *testing.T) {
		pq := NewDEPQ(3, 1)
		testlib.AclAssert(t, 1, pq.GetMin())
		testlib.AclAssert(t, 3, pq.GetMax())

		max := pq.PopMax()
		testlib.AclAssert(t, 3, max)
		testlib.AclAssert(t, 1, pq.Size())
		testlib.AclAssert(t, 1, pq.GetMin())
		testlib.AclAssert(t, 1, pq.GetMax())
	})
}

func BenchmarkDEPQ_PushPop(b *testing.B) {
	b.ReportAllocs()
	const m = 1_000_000
	rng := rand.New(rand.NewSource(1))
	vs := L1[int](m)
	for i := range vs {
		vs[i] = rng.Int()
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pq := NewDEPQ[int]()
		for _, v := range vs {
			pq.Push(v)
		}
		for !pq.Empty() {
			_ = pq.PopMin()
			if !pq.Empty() {
				_ = pq.PopMax()
			}
		}
	}
}
