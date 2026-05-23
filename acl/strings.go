package acl

// Manacher は文字が与えられたとき、各iについて、
// 文字iを中心とした回文の半径を記録した配列を返す
// 例）"ababa" => [0, 1, 2, 1, 0]
// O(|S|)
// 偶数調の回文を考慮する場合は、"a$b$a$b$a"のように$を挿入すると検出できるようになる
func Manacher(s string) []int {
	m := len(s)
	rad := make([]int, m)

	i, j := 0, 0
	for i < m {
		for i-j >= 0 && i+j < m && s[i-j] == s[i+j] {
			j++
		}
		rad[i] = j
		k := 1
		for i-k >= 0 && k+rad[i-k] < j {
			rad[i+k] = rad[i-k]
			k++
		}
		i += k
		j -= k
	}
	return rad
}

// kmpPrefix
func kmpPrefix(p string) []int {
	pi := L1[int](len(p))

	j := 0
	for i := 1; i < len(p); i++ {
		for j > 0 && p[j] != p[i] {
			j = pi[j-1]
		}
		if p[j] == p[i] {
			j++
		}
		pi[i] = j
	}
	return pi
}
