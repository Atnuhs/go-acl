package acl

import (
	"fmt"
	"math/bits"
)

type BIT struct {
	n    int
	data []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, data: make([]int, n+1)}
}

func BuildBIT(a []int) *BIT {
	n := len(a)
	b := &BIT{n: n, data: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		b.data[i] += a[i-1]
		j := i + (i & -i)
		if j <= n {
			b.data[j] += b.data[i]
		}
	}
	return b
}

func (b *BIT) Add(i, x int) {
	i++
	for i <= b.n {
		b.data[i] += x
		i += i & -i
	}
}

// Prefixは[0, .., i]の和
func (b *BIT) Prefix(i int) int {
	i++
	if i > b.n {
		panic(fmt.Errorf("BIT.Prefix out of range r should be <= %d but got %d", b.n, i))
	}
	ret := 0
	for i > 0 {
		ret += b.data[i]
		i -= i & -i
	}
	return ret
}

// Rangeは[l, r)の和
func (b *BIT) Range(l, r int) int {
	return b.Prefix(r-1) - b.Prefix(l-1)
}

// Atは一点取得
func (b *BIT) At(i int) int {
	return b.Prefix(i) - b.Prefix(i-1)
}

// SetはA[i]をvalに更新
func (b *BIT) Set(i, val int) {
	d := val - b.At(i)
	if d != 0 {
		b.Add(i, d)
	}
}

// FirstTrueAtはf(a[i])がTrueとなるような最小のiを返す
func (b *BIT) FirstTrueAt(f func(v int) bool) int {
	idx, sum := 0, 0
	step := 1 << (bits.Len(uint(b.n)) - 1)
	for step > 0 {
		next := idx + step
		if next <= b.n && !f(sum+b.data[next]) {
			sum += b.data[next]
			idx = next
		}
		step >>= 1
	}
	return idx
}
