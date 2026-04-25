package acl

import "math"

// Sieve はiが素数かをboolで持つテーブル
type Sieve []bool

// NewSieve はエラトステネスの篩の実装でテーブルを生成する
func NewSieve(n int) Sieve {
	isPrime := L1[bool](n + 1)
	for i := range isPrime {
		isPrime[i] = true
	}

	isPrime[0] = false
	isPrime[1] = false
	// sieve
	for i := 2; i*i <= n; i++ {
		if !isPrime[i] {
			continue
		}
		for j := i * i; j <= n; j += i {
			isPrime[j] = false
		}
	}
	return Sieve(isPrime)
}

func (sv Sieve) IsPrime(x int) bool {
	return sv[x]
}

func (sv Sieve) Primes() []int {
	n := len(sv) - 1
	m := func() int {
		if n < 17 {
			return [...]int{0, 0, 1, 2, 2, 3, 3, 4, 4, 4, 4, 5, 5, 6, 6, 6, 6}[n]
		}
		return int(1.3*float64(n)/math.Log(float64(n))) + 1
	}()
	ret := make([]int, 0, m)
	for i := 2; i <= n; i++ {
		if sv[i] {
			ret = append(ret, i)
		}
	}
	return ret
}

// MinFactorTable はある値xについてその最小の素因数を返すテーブル
type MinFactorTable []int

// NewMinFactor はO(N loglog N)でエラトステネスの篩の仕組みでminFactorテーブルを生成する
func NewMinFactor(n int) MinFactorTable {
	mf := make(MinFactorTable, n+1)
	mf[1] = 1
	for p := 2; p <= n; p++ {
		if mf[p] == 0 {
			mf[p] = p
			// p*pがオーバーフローする場合…あんまなさそう？
			if p > n/p {
				continue
			}
			for j := p * p; j <= n; j += p {
				if mf[j] == 0 {
					mf[j] = p
				}
			}
		}
	}
	return mf

}

// IsPrime はO(1)で素数かどうかを判定する
func (mf MinFactorTable) IsPrime(x int) bool {
	if x == 1 {
		return false
	}
	return mf[x] == x
}

// Factorize は O(log x)で素因数分解を行う
// 返り値は素因数とその指数のPairのスライス
// 例）got, ret
// 6, []Pair{{2,1}, {3.1}}
func (mf MinFactorTable) Factorize(x int) []Pair[int, int] {
	ret := make([]Pair[int, int], 0)
	n := x
	for n > 1 {
		p := mf[n]
		exp := 0

		for mf[n] == p {
			n /= p
			exp++
		}
		ret = append(ret, NewPair(p, exp))
	}
	return ret
}

// Divisors is O(log x) returns
// 2 => 1, 2
// 10 => 1, 2, 5, 10
func (mf *MinFactorTable) Divisors(x int) []int {
	ret := []int{1}

	f := mf.Factorize(x)
	for _, pe := range f {
		n := len(ret)
		for i := 0; i < n; i++ {
			v := 1
			for j := 0; j < pe.V; j++ {
				v *= pe.U
				ret = append(ret, ret[i]*v)
			}
		}
	}
	return ret
}

// CountDivisors is O(log x + len(sv.Divisors(x))) returns len(sv.Divisors(x))
// 1 => 1
// 2 => 2
// 10 => 4
func (mf *MinFactorTable) CountDivisors(x int) int {
	return CountDivisors(mf.Factorize(x))
}

// MobiusTable はメビウス関数μ(x)の値を持つテーブル
// メビウス関数は、整数nに対して以下のように定義される
// 0 <= n: nが平方数で割り切れる場合
// 1 or -1 <= (-1)^k: nがk個の異なる素因数を持つ場合
// 具体的には以下のような値となる
// 0 <= 4, 12, 18, 50: 平方数で割り切れる
// 1 <= 1, 6, 210: 偶数個の素因数を持つ
// -1 <= 2, 30, 140729 : 奇数個の素因数を持つ
// 約数系包除原理で使う
type MobiusTable []int

func NewMobiusTable(n int) MobiusTable {
	mf := NewMinFactor(n)
	mu := L1[int](n + 1)
	for i := range mu {
		mu[i] = 1
	}

	mu[1] = 1
	for x := 2; x <= n; x++ {
		p := mf[x]
		y := x / p
		if y%p == 0 {
			mu[x] = 0
		} else {
			mu[x] = -mu[y]
		}
	}
	return MobiusTable(mu)
}

func (mu MobiusTable) Mobius(x int) int {
	return mu[x]
}

// SegmentedSieveは
type SegmentedSieve struct {
	start   int
	isPrime Sieve
}

// 幅の狭い区間[L, R]に対して、D=R-Lとしたとき、
// 構築 O(Sqrt(R) log log R + D loglog R)で篩を生成することができる
func NewSegmentedSieve(lo, hi int) *SegmentedSieve {
	if lo >= hi {
		return nil
	}
	sqrtHi := Sqrt(hi) + 1
	sv1 := NewSieve(sqrtHi)

	m := hi - lo + 1
	isPrime := L1[bool](m)
	for i := range isPrime {
		isPrime[i] = true
	}

	for p := range sv1 {
		if !sv1.IsPrime(p) {
			continue
		}

		// l は lo 以上の最小のpの倍数
		l := ((lo + p - 1) / p) * p
		if l == p {
			l = p * p
		}
		for i := l; i <= hi; i += p {
			isPrime[i-lo] = false
		}
	}
	return &SegmentedSieve{
		start:   lo,
		isPrime: isPrime,
	}
}

func (sv *SegmentedSieve) IsPrime(x int) bool {
	i := x - sv.start
	return sv.isPrime[i]
}
