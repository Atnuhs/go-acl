package main

func main() {
	defer Out.Flush()

	n := 9
	r := Rs(n)

	ok := func(p Pos) bool {
		return InGrid(p, n, n) && GridAt(p, r) == '#'
	}

	ans := 0

	for i := range r {
		for j := range r[i] {
			p1 := NewPos(i, j)
			if !ok(p1) {
				continue
			}
			for di := -9; di <= 9; di++ {
				for dj := -9; dj <= 9; dj++ {
					if di == 0 && dj == 0 {
						continue
					}
					if !ok(p1.Add(di, dj)) {
						continue
					}
					if !ok(p1.Add(di-dj, dj+di)) {
						continue
					}
					if !ok(p1.Add(-dj, di)) {
						continue
					}

					ans++
				}
			}
		}
	}
	Ans(ans / 4)
}
