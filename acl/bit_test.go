package acl

import "fmt"

// 基本的な使い方: 一点加算と区間和クエリ。
func ExampleBIT() {
	bit := NewBIT(5) // 長さ 5 の BIT (全要素 0)
	bit.Add(0, 1)    // a = [1, 0, 0, 0, 0]
	bit.Add(2, 3)    // a = [1, 0, 3, 0, 0]
	bit.Add(4, 5)    // a = [1, 0, 3, 0, 5]

	fmt.Println(bit.Sum(0, 5)) // 区間 [0, 5) の和
	fmt.Println(bit.Sum(1, 4)) // 区間 [1, 4) の和
	fmt.Println(bit.At(2))     // a[2]
	// Output:
	// 9
	// 3
	// 3
}

// 既存スライスから BIT を構築する。
func ExampleBuildBIT() {
	bit := BuildBIT([]int{1, 2, 3, 4, 5})
	fmt.Println(bit.Sum(0, 5))
	fmt.Println(bit.Sum(1, 4))
	// Output:
	// 15
	// 9
}

// MaxRight は累積和の上で二分探索を行い、
// f(Sum(0, r)) が true となる最大の r を O(log n) で求める。
// 以下では「累積和がはじめて 10 以上になる手前の位置」を求めている。
func ExampleBIT_MaxRight() {
	bit := BuildBIT([]int{1, 2, 3, 4, 5}) // 累積和: 1, 3, 6, 10, 15
	r := bit.MaxRight(func(s int) bool { return s < 10 })
	fmt.Println(r)
	// Output:
	// 3
}
