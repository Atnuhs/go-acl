package acl

import (
	"reflect"
	"sort"
	"testing"
)

var primeTableTestCase = []struct {
	desc string
	x    int
	want bool
}{
	{desc: "one", x: 1, want: false},
	{desc: "small prime", x: 2, want: true},
	{desc: "small not prime", x: 4, want: false},
	{desc: "prime?", x: 57, want: false},
	{desc: "large prime", x: 104729, want: true},
	{desc: "large not prime", x: 111111, want: false},
}

func TestSieve(t *testing.T) {
	maxX := 3 * 100000
	sv := NewSieve(maxX)
	t.Log(sv[:100])

	for _, tc := range primeTableTestCase {
		t.Run(tc.desc, func(t *testing.T) {
			got := sv.IsPrime(tc.x)
			if got != tc.want {
				t.Errorf("%d is Prime?, expected %t, but got %t", tc.x, tc.want, got)
			}
		})
	}
}

func TestSieve_Primes(t *testing.T) {
	maxX := 3 * 100000
	sv := NewSieve(maxX)
	ps := sv.Primes()

	csv := 0
	for _, v := range sv {
		if v {
			csv++
		}
	}

	if csv != len(ps) {
		t.Errorf("Count Primes want: %d, got: %d", csv, len(ps))
	}

	for _, p := range ps {
		if !sv.IsPrime(p) {
			t.Errorf("%d should be prime but %t", p, sv[p])
		}
	}
}

func TestMinFactor_IsPrime(t *testing.T) {
	maxX := 3 * 100000
	mf := NewMinFactor(maxX)
	t.Log(mf[:100])

	for _, tc := range primeTableTestCase {
		t.Run(tc.desc, func(t *testing.T) {
			got := mf.IsPrime(tc.x)
			if got != tc.want {
				t.Errorf("%d is Prime?, expected %t, but got %t", tc.x, tc.want, got)
			}
		})
	}
}

func FuzzMinFactor_IsPrime(f *testing.F) {
	maxX := 3 * 100000
	mf := NewMinFactor(maxX)
	sv := NewSieve(maxX)
	f.Add(0)
	f.Add(1)
	f.Add(9)
	f.Add(123456)
	f.Add(1000000007)
	f.Fuzz(func(t *testing.T, a int) {
		if 1 > a || a > maxX {
			return
		}

		// Perform prime number determination in two different ways
		// method1 O(sqrt(N))
		ret1 := IsPrime(a)
		// method2 MinFactor
		ret2 := mf.IsPrime(a)
		// method3 Eratosthenes
		ret3 := sv.IsPrime(a)
		if ret1 != ret2 {
			t.Errorf("%d is Prime?, 試し割り: %t, but MinFactor %t", a, ret1, ret2)
		}
		if ret1 != ret3 {
			t.Errorf("%d is Prime?, 試し割り: %t, but Sieve %t", a, ret1, ret3)
		}
	})
}

func TestMinFactor_Factorize(t *testing.T) {
	maxX := 3 * 100000
	testCases := []struct {
		desc string
		x    int
		want []Pair[int, int]
	}{
		{desc: "one", x: 1, want: []Pair[int, int]{}},
		{desc: "simple prime number", x: 2, want: []Pair[int, int]{NewPair(2, 1)}},
		{desc: "simple composite number", x: 12, want: []Pair[int, int]{NewPair(2, 2), NewPair(3, 1)}},
		{desc: "large prime number", x: 104729, want: []Pair[int, int]{NewPair(104729, 1)}},
		{
			desc: "large composite number",
			x:    1260,
			want: []Pair[int, int]{NewPair(2, 2), NewPair(3, 2), NewPair(5, 1), NewPair(7, 1)},
		},
	}

	mf := NewMinFactor(maxX)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := mf.Factorize(tc.x)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("factorize %d result, expected %v, but got %v", tc.x, tc.want, got)
			}
		})
	}
}

func FuzzMinFactor_Factorize(f *testing.F) {
	maxX := 3 * 100000
	mf := NewMinFactor(maxX)
	f.Add(0)
	f.Add(1)
	f.Add(9)
	f.Add(123456)
	f.Add(1000000007)
	f.Fuzz(func(t *testing.T, a int) {
		if 1 > a || a > maxX {
			return
		}
		// Perform Factorization in two different ways
		// method1 O(sqrt(N))
		ret1 := Factorize(a)
		// method2 EratosthenesSieve
		ret2 := mf.Factorize(a)
		if !reflect.DeepEqual(ret1, ret2) {
			t.Errorf("factorize %d, method1: %v, but method2 %v", a, ret1, ret2)
		}
	})
}

func TestMinFactor_Divisors(t *testing.T) {
	maxX := 3 * 100000
	testCases := []struct {
		desc string
		x    int
		want []int
	}{
		{desc: "one", x: 1, want: []int{1}},
		{desc: "simple prime number", x: 2, want: []int{1, 2}},
		{desc: "simple composite number", x: 12, want: []int{1, 2, 3, 4, 6, 12}},
		{desc: "large prime number", x: 104729, want: []int{1, 104729}},
	}

	mf := NewMinFactor(maxX)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := mf.Divisors(tc.x)
			sort.Ints(got)
			sort.Ints(tc.want)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("factorize %d result, expected %v, but got %v", tc.x, tc.want, got)
			}
		})
	}
}

func FuzzMinFactor_Divisors(f *testing.F) {
	maxX := 3 * 100000
	mf := NewMinFactor(maxX)
	f.Add(0)
	f.Add(1)
	f.Add(9)
	f.Add(123456)
	f.Add(1000000007)
	f.Fuzz(func(t *testing.T, a int) {
		if 1 > a || a > maxX {
			return
		}
		// Enumerate divisors in two different ways
		// method1 O(sqrt(N))
		ret1 := Divisors(a)
		// method2 EratosthenesSieve
		ret2 := mf.Divisors(a)
		sort.Ints(ret1)
		sort.Ints(ret2)
		if !reflect.DeepEqual(ret1, ret2) {
			t.Errorf("enumerate %d divisors, method1: %v, but method2 %v", a, ret1, ret2)
		}
	})
}

func TestCountDivisors(t *testing.T) {
	maxX := 3 * 100000
	testCases := []struct {
		desc string
		x    int
		want int
	}{
		{desc: "one", x: 1, want: 1},
		{desc: "simple prime number", x: 2, want: 2},
		{desc: "simple composite number", x: 12, want: 6},
		{desc: "large prime number", x: 104729, want: 2},
	}

	mf := NewMinFactor(maxX)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			f := mf.Factorize(tc.x)
			got := CountDivisors(f)

			if tc.want != got {
				t.Errorf("%d divisors num, expected %d, but got %d", tc.x, tc.want, got)
			}
		})
	}
}

func FuzzCountDivisors(f *testing.F) {
	maxX := 3 * 100000
	mf := NewMinFactor(maxX)
	f.Add(0)
	f.Add(1)
	f.Add(9)
	f.Add(123456)
	f.Add(1000000007)
	f.Fuzz(func(t *testing.T, a int) {
		if 1 > a || a > maxX {
			return
		}
		// enumerate divisors counting methods
		type methodFunc func(x int) int
		methods := []methodFunc{
			// get len enumerate divisors
			func(x int) int {
				return len(Divisors(x))
			},
			// get len enumerate divisors with eratosthenes sieve
			func(x int) int {
				return len(mf.Divisors(x))
			},
			// get len from Factiroze
			func(x int) int {
				return CountDivisors(Factorize(x))
			},
			// get len from Factorize with sieve
			mf.CountDivisors,
		}

		for i, method1 := range methods {
			got1 := method1(a)
			for j := i + 1; j < len(methods); j++ {
				method2 := methods[j]
				got2 := method2(a)
				t.Logf(
					"%d num divisors, method1: %d, but method2 %d",
					a,
					got1,
					got2,
				)
				if !reflect.DeepEqual(got1, got2) {
					t.Errorf(
						"%d num divisors, method1: %v:%v, but method2 %v:%v",
						a,
						method1,
						got1,
						method2,
						got2,
					)
				}
			}
		}
	})
}

func TestMobius(t *testing.T) {
	maxX := 3 * 100000
	testCases := []struct {
		desc string
		x    int
		want int
	}{
		{desc: "one", x: 1, want: 1},
		{desc: "simple prime number", x: 2, want: -1},
		{desc: "simple composite number", x: 6, want: 1},
		{desc: "simple composite number", x: 30, want: -1},
		{desc: "simple composite number", x: 4, want: 0},
		{desc: "simple composite number", x: 12, want: 0},
		{desc: "large prime number", x: 104729, want: -1},
	}

	mu := NewMobiusTable(maxX)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := mu.Mobius(tc.x)

			if tc.want != got {
				t.Errorf("%d Mebius function, expected %d, but got %d", tc.x, tc.want, got)
			}
		})
	}
}

func FuzzMobius(f *testing.F) {
	maxX := 3 * 100000
	mu := NewMobiusTable(maxX)
	f.Add(0)
	f.Add(1)
	f.Add(9)
	f.Add(123456)
	f.Add(1000000007)
	f.Fuzz(func(t *testing.T, a int) {
		if 1 > a || a > maxX {
			return
		}
		// Mobius function in two different ways
		// method1 O(sqrt(N))
		ret1 := Mobius(a)
		// method2 EratosthenesSieve
		ret2 := mu.Mobius(a)
		if ret1 != ret2 {
			t.Errorf("%d mobius function, method1: %v, but method2 %v", a, ret1, ret2)
		}
	})
}

func TestSegmentedSieve(t *testing.T) {
	testcase := map[string]struct {
		lo, hi int
	}{
		"small":  {lo: 2, hi: 10},
		"medium": {lo: 1000000, hi: 2000000},
		"large":  {lo: 10_000_000_000, hi: 10_000_100_000},
	}

	for name, tc := range testcase {
		t.Run(name, func(t *testing.T) {
			sv := NewSegmentedSieve(tc.lo, tc.hi)

			for i := tc.lo; i <= tc.hi; i++ {
				got := sv.IsPrime(i)
				want := IsPrime(i)

				if got != want {
					t.Errorf("%d is Prime? => want: %t, got: %t", i, want, got)
				}
			}
		})
	}
}

func FuzzSegmentedSieve(f *testing.F) {
	MaxHi := 1_000_000_000_000
	MaxL := 1_000_000

	f.Add(2, 10)
	f.Add(1000, 2000)
	f.Add(10_000_000_000, 10_000_100_000)
	f.Fuzz(func(t *testing.T, lo, hi int) {
		if lo <= 2 || hi < lo {
			return
		}
		if hi > MaxHi {
			return
		}
		l := hi - lo
		if l < 1 || l > MaxL {
			return
		}

		ssv := NewSegmentedSieve(lo, hi)
		for i := lo; i <= hi; i++ {
			got := ssv.IsPrime(i)
			want := IsPrime(i)

			if got != want {
				t.Errorf("[lo=%d,hi=%d] %d is Prime? => want: %t, got: %t", lo, hi, i, want, got)
			}
		}
	})
}
