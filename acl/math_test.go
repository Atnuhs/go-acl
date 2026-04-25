package acl

import (
	"sort"
	"testing"
)

func TestModPow(t *testing.T) {
	tests := []struct {
		x, e, mod, want int
	}{
		{2, 10, 1000000007, 1024},
		{3, 0, 7, 1},
		{5, 1, 7, 5},
		{2, 3, 5, 3},    // 8 % 5 = 3
		{0, 5, 7, 0},
		{7, 2, 1, 0},    // mod <= 1
		{3, -1, 7, 0},   // 負の指数
	}
	for _, tc := range tests {
		got := ModPow(tc.x, tc.e, tc.mod)
		if got != tc.want {
			t.Errorf("ModPow(%d,%d,%d) = %d, want %d", tc.x, tc.e, tc.mod, got, tc.want)
		}
	}
}

func TestLcm(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{4, 6, 12},
		{3, 5, 15},
		{12, 18, 36},
		{1, 7, 7},
	}
	for _, tc := range tests {
		got := Lcm(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("Lcm(%d,%d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		x, want int
	}{
		{0, 0},
		{1, 1},
		{4, 2},
		{9, 3},
		{2, 1},
		{8, 2},
		{100, 10},
		{99, 9},
	}
	for _, tc := range tests {
		got := Sqrt(tc.x)
		if got != tc.want {
			t.Errorf("Sqrt(%d) = %d, want %d", tc.x, got, tc.want)
		}
	}
}

func TestSqrt_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for negative input")
		}
	}()
	Sqrt(-1)
}

func TestNextPerm(t *testing.T) {
	a := []int{1, 2, 3}
	perms := [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}}
	for i, want := range perms {
		for j := range a {
			if a[j] != want[j] {
				t.Errorf("perm %d: got %v, want %v", i, a, want)
				break
			}
		}
		NextPerm(a)
	}

	// 最後の並び替えではfalseを返す
	last := []int{3, 2, 1}
	if NextPerm(last) {
		t.Error("NextPerm on last permutation should return false")
	}
}

func TestExtrema(t *testing.T) {
	mi, ma := Extrema(3, 1, 4, 1, 5, 9, 2, 6)
	if mi != 1 || ma != 9 {
		t.Errorf("Extrema = (%d, %d), want (1, 9)", mi, ma)
	}
	mi2, ma2 := Extrema(42)
	if mi2 != 42 || ma2 != 42 {
		t.Errorf("Extrema(42) = (%d, %d), want (42, 42)", mi2, ma2)
	}
}

func TestMax(t *testing.T) {
	if got := Max(3, 1, 4, 1, 5); got != 5 {
		t.Errorf("Max = %d, want 5", got)
	}
	if got := Max(-3, -1); got != -1 {
		t.Errorf("Max(-3,-1) = %d, want -1", got)
	}
}

func TestMin(t *testing.T) {
	if got := Min(3, 1, 4, 1, 5); got != 1 {
		t.Errorf("Min = %d, want 1", got)
	}
	if got := Min(-3, -1); got != -3 {
		t.Errorf("Min(-3,-1) = %d, want -3", got)
	}
}

func TestSum(t *testing.T) {
	if got := Sum(1, 2, 3, 4, 5); got != 15 {
		t.Errorf("Sum = %d, want 15", got)
	}
	if got := Sum(-1, 1); got != 0 {
		t.Errorf("Sum(-1,1) = %d, want 0", got)
	}
}

func TestAbs(t *testing.T) {
	tests := []struct{ x, want int }{{5, 5}, {-5, 5}, {0, 0}}
	for _, tc := range tests {
		if got := Abs(tc.x); got != tc.want {
			t.Errorf("Abs(%d) = %d, want %d", tc.x, got, tc.want)
		}
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		x    int
		want bool
	}{
		{1, false}, {2, true}, {3, true}, {4, false},
		{17, true}, {25, false}, {97, true}, {100, false},
	}
	for _, tc := range tests {
		if got := IsPrime(tc.x); got != tc.want {
			t.Errorf("IsPrime(%d) = %v, want %v", tc.x, got, tc.want)
		}
	}
}

func TestFactorize(t *testing.T) {
	tests := []struct {
		x    int
		want []Pair[int, int]
	}{
		{1, []Pair[int, int]{}},
		{2, []Pair[int, int]{{2, 1}}},
		{4, []Pair[int, int]{{2, 2}}},
		{12, []Pair[int, int]{{2, 2}, {3, 1}}},
		{30, []Pair[int, int]{{2, 1}, {3, 1}, {5, 1}}},
	}
	for _, tc := range tests {
		got := Factorize(tc.x)
		if len(got) != len(tc.want) {
			t.Errorf("Factorize(%d) = %v, want %v", tc.x, got, tc.want)
			continue
		}
		for i := range got {
			if got[i].U != tc.want[i].U || got[i].V != tc.want[i].V {
				t.Errorf("Factorize(%d)[%d] = %v, want %v", tc.x, i, got[i], tc.want[i])
			}
		}
	}
}

func TestDivisors(t *testing.T) {
	tests := []struct {
		x    int
		want []int
	}{
		{1, []int{1}},
		{2, []int{1, 2}},
		{10, []int{1, 2, 5, 10}},
		{12, []int{1, 2, 3, 4, 6, 12}},
	}
	for _, tc := range tests {
		got := Divisors(tc.x)
		sort.Ints(got)
		sort.Ints(tc.want)
		if len(got) != len(tc.want) {
			t.Errorf("Divisors(%d) = %v, want %v", tc.x, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("Divisors(%d) = %v, want %v", tc.x, got, tc.want)
				break
			}
		}
	}
}

func TestCeilDiv(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{10, 3, 4},
		{9, 3, 3},
		{0, 5, 0},
		{1, 2, 1},
	}
	for _, tc := range tests {
		if got := CeilDiv(tc.a, tc.b); got != tc.want {
			t.Errorf("CeilDiv(%d,%d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestCeilPow2(t *testing.T) {
	tests := []struct{ x, want int }{
		{1, 1}, {2, 2}, {3, 4}, {4, 4}, {5, 8}, {8, 8},
	}
	for _, tc := range tests {
		if got := CeilPow2(tc.x); got != tc.want {
			t.Errorf("CeilPow2(%d) = %d, want %d", tc.x, got, tc.want)
		}
	}
}

func TestPow(t *testing.T) {
	tests := []struct{ x, e, want int }{
		{2, 0, 1}, {2, 10, 1024}, {3, 3, 27}, {5, 4, 625},
	}
	for _, tc := range tests {
		if got := Pow(tc.x, tc.e); got != tc.want {
			t.Errorf("Pow(%d,%d) = %d, want %d", tc.x, tc.e, got, tc.want)
		}
	}
}

func TestPow10(t *testing.T) {
	tests := []struct{ e, want int }{
		{0, 1}, {1, 10}, {3, 1000}, {6, 1000000},
	}
	for _, tc := range tests {
		if got := Pow10(tc.e); got != tc.want {
			t.Errorf("Pow10(%d) = %d, want %d", tc.e, got, tc.want)
		}
	}
}

func TestInv(t *testing.T) {
	testCases := []struct {
		desc string
		x    int
		p    int
	}{
		{desc: "x:20, p:7", x: 20, p: 7},
		{desc: "x:1234567, p:10^9+7", x: 1234567, p: 1000000007},
		{desc: "x:1234567, p:998244353", x: 1234567, p: 998244353},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			invX := Inv(tc.x, tc.p)

			// x * invX = should be 1
			got := (tc.x * invX) % tc.p
			if got != 1 {
				t.Errorf("actual should be 1 but got %d, invX: %d", got, invX)
			}
		})
	}
}

func FuzzInv(f *testing.F) {
	f.Add(4, 7)
	f.Add(1000000007, 97)
	f.Add(97, 1000000007)
	f.Add(4, 1)
	f.Add(4, -11)
	f.Fuzz(func(f *testing.T, x, mod int) {
		if !IsPrime(mod) || mod <= 1 {
			return
		}

		if Gcd(x, mod) != 1 {
			return
		}

		invX := Inv(x, mod)
		got := (invX * x) % mod

		if got != 1 {
			f.Errorf("expected 1, but got %d, x: %d, mod: %d, invX: %d", got, x, mod, invX)
		}
	})
}

func TestGcd(t *testing.T) {
	testCases := []struct {
		desc       string
		x, y, want int
	}{
		{desc: "gcd(2, 2) => 2", x: 2, y: 2, want: 2},
		{desc: "gcd(4, 2) => 2", x: 4, y: 2, want: 2},
		{desc: "gcd(4, 6) => 2", x: 4, y: 6, want: 2},
		{desc: "gcd(11, 13) => 1", x: 11, y: 13, want: 1},
		{desc: "gcd(11, 13) => 1", x: 11, y: 13, want: 1},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := Gcd(tc.x, tc.y)
			if tc.want != got {
				t.Errorf("expected %d but got %d", tc.want, got)
			}
		})
	}
}

func FuzzGcd(f *testing.F) {
	f.Add(0, 100000)
	f.Add(100000, 0)
	f.Add(0, 0)
	f.Add(12345678, 1000000007)
	f.Add(-12345678, 1000000007)
	f.Add(-12345678, 0)
	f.Add(0, -12345678)
	f.Fuzz(func(f *testing.T, x, y int) {
		_ = Gcd(x, y)
	})
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		desc     string
		x        int
		mod      int
		expected int
	}{
		{desc: "factorial 0", x: 0, mod: 7, expected: 1},
		{desc: "factorial 1", x: 1, mod: 7, expected: 1},
		{desc: "factorial 3", x: 3, mod: 7, expected: 6},
		{desc: "factorial 5", x: 5, mod: 7, expected: 1}, // 5! = 120, 120 % 7 = 1
		{desc: "negative x", x: -1, mod: 7, expected: 0},
		{desc: "mod <= 1", x: 5, mod: 1, expected: 0},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := Factorial(tc.x, tc.mod)
			if got != tc.expected {
				t.Errorf("expected %d but got %d", tc.expected, got)
			}
		})
	}
}
