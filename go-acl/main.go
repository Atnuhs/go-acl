package main

func main() {
	defer Out.Flush()

	n := I()
	as, bs := IIs(n)

	m := I()
	ss := Bs(m)

	mp := L3[int](10, 10, 27)
	for _, s := range ss {
		for i, r := range s {
			mp[len(s)-1][i][r-'a']++
		}
	}

	for _, s := range ss {
		YesNoFunc(func() bool {
			if len(s) != n {
				return false
			}
			for j, r := range s {
				a, b := as[j]-1, bs[j]-1
				if mp[a][b][r-'a'] == 0 {
					return false
				}
			}
			return true
		})
	}
}
