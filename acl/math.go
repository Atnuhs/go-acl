package acl

import (
	"cmp"
	"math"
)

// return x! % mod with O(x)
func Factorial(x, mod int) int {
	if x < 0 || mod <= 1 {
		return 0
	}
	if x == 0 || x == 1 {
		return 1
	}
	ans := 1
	for i := 2; i <= x; i++ {
		ans = (ans * i) % mod
	}
	return ans
}

// ModFact precomputes n!, (n!)^-1 mod p in O(n). mod must be prime and > n.
type ModFact struct {
	mod   int
	fact  []int
	ifact []int
}

// NewModFact precomputes factorials up to n. mod must be prime and > n.
func NewModFact(n, mod int) *ModFact {
	fact := make([]int, n+1)
	ifact := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	ifact[n] = ModPow(fact[n], mod-2, mod)
	for i := n; i > 0; i-- {
		ifact[i-1] = ifact[i] * i % mod
	}
	return &ModFact{mod: mod, fact: fact, ifact: ifact}
}

// 【n!】Fact returns n! mod p. O(1)。n は前計算サイズ以下。
func (f *ModFact) Fact(n int) int {
	if n < 0 || n >= len(f.fact) {
		return 0
	}
	return f.fact[n]
}

// 【(n!)^-1】IFact returns (n!)^-1 mod p. O(1)。n は前計算サイズ以下。
func (f *ModFact) IFact(n int) int {
	if n < 0 || n >= len(f.ifact) {
		return 0
	}
	return f.ifact[n]
}

// 【nPr / 順列】Perm returns nPr mod p. O(1)。
// n, r ともに前計算サイズ以下のときに使う。
func (f *ModFact) Perm(n, r int) int {
	if r < 0 || n < r || n >= len(f.fact) {
		return 0
	}
	return f.fact[n] * f.ifact[n-r] % f.mod
}

// 【nCr / 組合せ・通常版】Comb returns nCr mod p. O(1)。
// n が前計算サイズ以下のときに使う（例: n <= 2*10^6）。
// n が巨大なら CombBigN を使うこと。
func (f *ModFact) Comb(n, r int) int {
	if r < 0 || n < r || n >= len(f.fact) {
		return 0
	}
	return f.fact[n] * f.ifact[r] % f.mod * f.ifact[n-r] % f.mod
}

// 【nCr / 組合せ・nが巨大版】CombBigN returns nCr mod p. O(r)。
// n が前計算サイズを超えるが r は小さいとき用（例: n <= 10^18, r <= 10^5）。
// 前計算は r 以上あればよい。n が小さいなら Comb の方が高速。
func (f *ModFact) CombBigN(n, r int) int {
	if r < 0 || n < r || r >= len(f.fact) {
		return 0
	}
	ret := 1
	for i := range r {
		ret = (ret * (n - i)) % f.mod
	}
	return (ret * f.ifact[r]) % f.mod
}

// Stirling2 returns S(m,n) mod p in O(N log M).
// ModFact must be precomputed for at least n (i.e. NewModFact(>=n, mod)).
func (f *ModFact) Stirling2(m, n int) int {
	if n >= len(f.fact) {
		panic("ModFact.Stirling2: n exceeds precomputed range")
	}
	// 包除原理を利用する
	ret := 0
	for k := 0; k <= n; k++ {
		term := (f.Comb(n, k) * ModPow(n-k, m, f.mod)) % f.mod
		if k&1 == 0 {
			ret += term
		} else {
			ret += f.mod - term
		}
		ret %= f.mod
	}
	return (ret * f.ifact[n]) % f.mod
}

// 【nCr / 組合せ・パスカル表版】CombTable returns table t where t[i][j] = iCj mod p (0<=i<=n, 0<=j<=r, j>i は 0)。
// O(N*R) で表を構築、N*R <= 10^6 が目安。ModFact 不要・mod が合成数でも可。
// 用途: 多くの (i,j) を引きたい / mod が素数でない / ModFact を作りたくないとき。
// 単発の nCr なら Comb / CombBigN の方が軽い。
func CombTable(n, r, mod int) [][]int {
	t := L2[int](n+1, r+1)
	for i := 0; i <= n; i++ {
		t[i][0] = 1
		upper := min(i, r)
		for j := 1; j <= upper; j++ {
			t[i][j] = (t[i-1][j-1] + t[i-1][j]) % mod
		}
	}
	return t
}

// ModPow return x^e % mod in O(log x).
func ModPow(x, e, mod int) int {
	if mod <= 1 {
		return 0
	}
	if e == 0 {
		return 1
	}
	if e < 0 {
		return 0 // 負の指数は扱わない
	}

	ret := 1
	x %= mod
	for e > 0 {
		if e&1 == 1 {
			ret = (ret * x) % mod
		}
		x = (x * x) % mod
		e >>= 1
	}
	return ret
}

// Inv return x^(-1) % mod
func Inv(x, mod int) int {
	return ModPow(x, mod-2, mod)
}

// Gcd return greatest common divisor on O(log N)
func Gcd(a, b int) int {
	// Handle negative numbers
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	// Iterative implementation to avoid stack overflow
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Lcm return least common multiple on O(log N)
func Lcm(a, b int) int {
	return a / Gcd(a, b) * b
}

// Sqrt return square root of x
func Sqrt(x int) int {
	if x < 0 {
		panic("Sqrt negative input")
	}
	if x <= 1 {
		return x
	}
	ret := int(math.Sqrt(float64(x)))
	for ret > x/ret {
		ret--
	}
	for (ret + 1) <= x/(ret+1) {
		ret++
	}
	return ret
}

// NextPerm returns [1,2,3,4] => [1,2,4,3] ... [4,3,2,1]
func NextPerm(a []int) bool {
	// search i
	i := len(a) - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := len(a) - 1
	for j >= 0 && a[j] <= a[i] {
		j--
	}

	a[i], a[j] = a[j], a[i]

	l := i + 1
	r := len(a) - 1
	for l < r {
		a[l], a[r] = a[r], a[l]
		l++
		r--
	}
	return true
}

// Extrema returns min, max
func Extrema[T cmp.Ordered](vals ...T) (T, T) {
	mi, ma := vals[0], vals[0]
	for _, v := range vals {
		if v < mi {
			mi = v
		}
		if v > ma {
			ma = v
		}
	}
	return mi, ma
}

func Max[T cmp.Ordered](vals ...T) T {
	ma := vals[0]
	for _, v := range vals[1:] {
		if v > ma {
			ma = v
		}
	}
	return ma
}

func Min[T cmp.Ordered](vals ...T) T {
	mi := vals[0]
	for _, v := range vals[1:] {
		if v < mi {
			mi = v
		}
	}
	return mi
}

// Sum returns sum of vals
func Sum[T cmp.Ordered](vals ...T) T {
	var sum T
	for _, v := range vals {
		sum += v
	}
	return sum
}

// Abs returns absolute value of x
func Abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

// IsPrime is O(Sqrt(N))
func IsPrime(x int) bool {
	if x <= 1 {
		return false
	}
	if x <= 3 {
		return true
	}
	if x%2 == 0 || x%3 == 0 {
		return false
	}

	rx := Sqrt(x)
	for i := 5; i <= rx; i += 6 {
		if x%i == 0 || x%(i+2) == 0 {
			return false
		}
	}
	return true
}

// Factorize is O(Sqrt(N))
// got, ret
// 6, []Pair{{2,1}, {3.1}}
func Factorize(x int) []Pair[int, int] {
	if x <= 1 {
		return []Pair[int, int]{}
	}

	ret := make([]Pair[int, int], 0)
	n := x

	// Handle factor 2
	if n%2 == 0 {
		exp := 0
		for n%2 == 0 {
			n /= 2
			exp++
		}
		ret = append(ret, NewPair(2, exp))
	}

	// Handle odd factors from 3 onwards
	rx := Sqrt(n)
	for i := 3; i <= rx; i += 2 {
		if n%i == 0 {
			exp := 0
			for n%i == 0 {
				n /= i
				exp++
			}
			ret = append(ret, NewPair(i, exp))
			rx = Sqrt(n) // Update rx after reducing n
		}
	}

	if n > 1 {
		ret = append(ret, NewPair(n, 1))
	}
	return ret
}

// Mobius is O(sqrt(n)) returns
// 0 <= 4, 12, 18, 50
// 1 <= 1, 6, 210
// -1 <= 2, 30, 140729
func Mobius(x int) int {
	if x <= 0 {
		return 0
	}
	if x == 1 {
		return 1
	}

	ret := 1
	n := x

	// Handle factor 2
	if n%2 == 0 {
		n /= 2
		ret = -ret
		if n%2 == 0 { // Square factor
			return 0
		}
	}

	// Handle odd factors
	rx := Sqrt(n)
	for i := 3; i <= rx; i += 2 {
		if n%i == 0 {
			n /= i
			ret = -ret
			if n%i == 0 { // Square factor
				return 0
			}
			rx = Sqrt(n)
		}
	}

	if n > 1 {
		ret = -ret
	}
	return ret
}

// Divisors is O(sqrt(n)) returns
// 2 => 1, 2
// 10 => 1, 2, 5, 10
func Divisors(x int) []int {
	if x <= 0 {
		return []int{}
	}

	ret := make([]int, 0)
	rx := Sqrt(x)
	for i := 1; i <= rx; i++ {
		if x%i == 0 {
			ret = append(ret, i)
			if i != x/i {
				ret = append(ret, x/i)
			}
		}
	}
	return ret
}

// CountDivisors is O(sqrt(n)) returns
// 1 => 1
// 2 => 2
// 10 => 4
func CountDivisors(pairs []Pair[int, int]) int {
	ans := 1
	for _, pe := range pairs {
		ans *= (pe.V + 1)
	}
	return ans
}

// CeilDiv は a / bの切り上げを返す
func CeilDiv(a, b int) int {
	if b <= 0 {
		panic("b must be > 0")
	}
	return (a + b - 1) / b
}

// CeilPow2はx以上の最小の2冪を返す
func CeilPow2(x int) int {
	ret := 1
	for ret < x {
		ret <<= 1
	}
	return ret
}

func Pow(x, e int) int {
	ret := 1
	for e > 0 {
		if e&1 == 1 {
			ret *= x
		}
		x *= x
		e >>= 1
	}
	return ret
}

func Pow10(e int) int {
	ret := 1
	for range e {
		ret *= 10
	}
	return ret
}

// Stirling2 はO(MN)で第2種スターリング数を計算する（MODなし）
func Stirling2(m, n int) (int, []int) {
	// 二次元DPを用いる（1次元化している）
	// S(M,N) = S(m-1,n-1) + n * S(m-1,n)

	dp := L1[int](n + 1)
	dp[0] = 1
	for range m {
		for j := n; j > 0; j-- {
			dp[j] = j*dp[j] + dp[j-1]
		}
		dp[0] = 0
	}
	return dp[n], dp
}

// Stirling2Mod はO(MN)で第2種スターリング数を計算する（MODあり）
//
//	M <= 10^6, N <= 10^6, MN <= 10^6まで対応。
//	M <= 10^18, N <= 10^5 のような制約の場合は、ModFact のStirling2を使うこと
func Stirling2Mod(m, n, mod int) (int, []int) {
	// 二次元DPを用いる（1次元化している）
	// S(M,N) = S(m-1,n-1) + n * S(m-1,n)

	dp := L1[int](n + 1)
	dp[0] = 1
	for range m {
		for j := n; j > 0; j-- {
			dp[j] = ((j*dp[j])%mod + dp[j-1]) % mod
		}
		dp[0] = 0
	}
	return dp[n], dp
}
