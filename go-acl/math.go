package main

import (
	"math"

	"golang.org/x/exp/constraints"
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

// ModPow return x^e % mod
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
func Extrema[T constraints.Ordered](vals ...T) (T, T) {
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

func Max[T constraints.Ordered](vals ...T) T {
	ma := vals[0]
	for _, v := range vals[1:] {
		if v > ma {
			ma = v
		}
	}
	return ma
}

func Min[T constraints.Ordered](vals ...T) T {
	mi := vals[0]
	for _, v := range vals[1:] {
		if v < mi {
			mi = v
		}
	}
	return mi
}

// Sum returns sum of vals
func Sum[T constraints.Ordered](vals ...T) T {
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
func Factorize(x int) []*Pair[int, int] {
	if x <= 1 {
		return []*Pair[int, int]{}
	}

	ret := make([]*Pair[int, int], 0)
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
func CountDivisors(pairs []*Pair[int, int]) int {
	ans := 1
	for _, pe := range pairs {
		ans *= (pe.V + 1)
	}
	return ans
}

// CeilPow2はx以上の最小の2冪を返す
func CeilPow2(x int) int {
	ret := 1
	for ret < x {
		ret <<= 1
	}
	return ret
}
