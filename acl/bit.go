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

// Sumは[l, r)の和を返す
func (b *BIT) Sum(l, r int) int {
	if l < 0 || r > b.n || l > r {
		panic(fmt.Errorf("BIT.Sum out of range: must be 0 <= l <= r <= %d but got l=%d, r=%d", b.n, l, r))
	}
	return b.prefix(r) - b.prefix(l)
}

// prefixは内部用。a[0]+...+a[i-1] を返す（i は 1-indexed の終端）
func (b *BIT) prefix(i int) int {
	ret := 0
	for i > 0 {
		ret += b.data[i]
		i -= i & -i
	}
	return ret
}

// Prefixは[0, i]の和を返す
func (b *BIT) Prefix(i int) int {
	return b.Sum(0, i+1)
}

// Atはa[i]を返す
func (b *BIT) At(i int) int {
	return b.Sum(i, i+1)
}

// SetはA[i]をvalに更新
func (b *BIT) Set(i, val int) {
	d := val - b.At(i)
	if d != 0 {
		b.Add(i, d)
	}
}

// MaxRightはf(Sum(0, r))がTrueとなるような最大のrを返す。
// f(0)がTrue、かつf は単調 (True → False) であることを要求する。
func (b *BIT) MaxRight(f func(v int) bool) int {
	if b.n == 0 {
		return 0
	}
	idx, sum := 0, 0
	step := 1 << (bits.Len(uint(b.n)) - 1)
	for step > 0 {
		next := idx + step
		if next <= b.n && f(sum+b.data[next]) {
			sum += b.data[next]
			idx = next
		}
		step >>= 1
	}
	return idx
}
