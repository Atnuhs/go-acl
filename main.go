package main

import . "github.com/Atnuhs/go-acl/acl"

func main() {
	defer Out.Flush()

	s, t := B(), B()
	n, m := len(s), len(t)

	next := L2[int](n+1, m)
	F2(next, -1)
	for i := n - 1; i >= 0; i-- {
		for j := range m {
			if s[i] == t[j] {
				next[i][j] = i
			} else {
				next[i][j] = next[i+1][j]
			}
		}
	}

	find := func(cur int) int {
		for i := range m {
			cur = next[cur][i] + 1
			if cur == 0 {
				return -1
			}
		}
		return cur
	}

	cnt := 0
	last := -1
	for i := range s {
		j := find(i)
		if j == -1 {
			Ans((n*(n+1))/2 - cnt)
			return
		}
		ta, tb := (i - last), (n - j + 1)
		cnt += ta * tb
		last = i
	}
	Ans((n*(n+1))/2 - cnt)
}
