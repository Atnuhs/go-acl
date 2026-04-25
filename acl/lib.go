package acl

import (
	"cmp"
	"fmt"
	"strconv"
)

const (
	MOD1 = 1000000007
	MOD2 = 998244353
	// INF is 10^18
	INF = 1000000000000000000

	// Error messages for data structures
	ErrEmptyContainer = "operation on empty container"
	ErrOutOfIndex     = "index out of range"
)

type Entry[K cmp.Ordered, V any] struct {
	K K
	V V
}

type Ok[T any] func(x T) bool

func S2i(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func I2s(i int) string {
	return fmt.Sprint(i)
}

func InRange(x, l, r int) bool {
	return l <= x && x < r
}
