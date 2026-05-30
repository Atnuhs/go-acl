package acl

// DEQUE_CAP は Deque の初期バッファ長 (2 のべき乗)。
// 容量不足になると 2 倍ずつ拡張される。
const DEQUE_CAP = 1 << 10

// Deque はリングバッファ実装による両端キュー (double-ended queue)。
// 前後への追加・削除がともに償却 O(1) で行える。
// 内部バッファ長は常に 2 のべき乗で、満杯になると倍々に拡張される。
type Deque[T any] struct {
	buf  []T
	l, r int
}

// NewDeque は空の Deque を生成する。初期容量は DEQUE_CAP。
//
// 計算量: O(DEQUE_CAP)
func NewDeque[T any]() *Deque[T] {
	buf := make([]T, DEQUE_CAP)
	return &Deque[T]{
		buf: buf,
	}
}

func (d *Deque[T]) next(r int) int { return (r + 1) & (len(d.buf) - 1) }
func (d *Deque[T]) prev(r int) int { return (r - 1) & (len(d.buf) - 1) }

func (d *Deque[T]) size() int {
	return (d.r - d.l + len(d.buf)) & (len(d.buf) - 1)
}

func (d *Deque[T]) isEmpty() bool {
	return d.l == d.r
}

func (d *Deque[T]) isFull() bool {
	return d.next(d.r) == d.l
}

func (d *Deque[T]) grow() {
	size := d.size()

	old := d.buf
	n := len(old) << 1
	d.buf = make([]T, n)

	// データをコピー
	if d.l <= d.r {
		// 連続している場合
		copy(d.buf, old[d.l:d.r])
	} else {
		// 循環している場合
		p := copy(d.buf, old[d.l:])
		copy(d.buf[p:], old[:d.r])
	}

	d.l = 0
	d.r = size
}

// GetFront は先頭要素を返す (削除しない)。
// 空のときの挙動は未定義のため、呼び出し側で IsEmpty を確認すること。
//
// 計算量: O(1)
func (d *Deque[T]) GetFront() T {
	return d.buf[d.l]
}

// GetBack は末尾要素を返す (削除しない)。
// 空のときの挙動は未定義のため、呼び出し側で IsEmpty を確認すること。
//
// 計算量: O(1)
func (d *Deque[T]) GetBack() T {
	return d.buf[d.prev(d.r)]
}

// PushFront は先頭に v を追加する。
//
// 計算量: 償却 O(1) (拡張時のみ O(n))
func (d *Deque[T]) PushFront(v T) {
	if d.isFull() {
		d.grow()
	}
	d.l = d.prev(d.l)
	d.buf[d.l] = v // ← 修正: 値を設定
}

// PushBack は末尾に v を追加する。
//
// 計算量: 償却 O(1) (拡張時のみ O(n))
func (d *Deque[T]) PushBack(v T) {
	if d.isFull() {
		d.grow()
	}
	d.buf[d.r] = v
	d.r = d.next(d.r)
}

// PopFront は先頭要素を取り出して返す。空のときは panic する。
//
// 計算量: O(1)
func (d *Deque[T]) PopFront() T {
	if d.isEmpty() {
		panic("Deque PopFront deque is empty")
	}
	v := d.buf[d.l]
	d.l = d.next(d.l)
	return v
}

// PopBack は末尾要素を取り出して返す。空のときは panic する。
//
// 計算量: O(1)
func (d *Deque[T]) PopBack() T {
	if d.isEmpty() {
		panic("Deque PopBack deque is empty")
	}
	d.r = d.prev(d.r)
	return d.buf[d.r]
}

// At は先頭から i 番目 (0-indexed) の要素を返す。
// i が [0, Size()) の範囲外なら panic する。
//
// 計算量: O(1)
func (d *Deque[T]) At(i int) T {
	if i < 0 || i >= d.size() {
		panic("deque: index out of range")
	}
	return d.buf[(d.l+i)&(len(d.buf)-1)]
}

// Update は先頭から i 番目 (0-indexed) の要素を v で上書きする。
// i が [0, Size()) の範囲外なら panic する。
//
// 計算量: O(1)
func (d *Deque[T]) Update(i int, v T) {
	if i < 0 || i >= d.size() {
		panic("deque: index out of range")
	}
	d.buf[(d.l+i)&(len(d.buf)-1)] = v
}

// Size は現在の要素数を返す。
//
// 計算量: O(1)
func (d *Deque[T]) Size() int {
	return d.size()
}

// IsEmpty は要素数が 0 のとき true を返す。
//
// 計算量: O(1)
func (d *Deque[T]) IsEmpty() bool {
	return d.isEmpty()
}
