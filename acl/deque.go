package acl

const DEQUE_CAP = 1 << 10

type Deque[T any] struct {
	buf  []T
	l, r int
}

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

func (d *Deque[T]) GetFront() T {
	return d.buf[d.l]
}

func (d *Deque[T]) GetBack() T {
	return d.buf[d.r]
}

func (d *Deque[T]) PushFront(v T) {
	if d.isFull() {
		d.grow()
	}
	d.l = d.prev(d.l)
	d.buf[d.l] = v // ← 修正: 値を設定
}

func (d *Deque[T]) PushBack(v T) {
	if d.isFull() {
		d.grow()
	}
	d.buf[d.r] = v
	d.r = d.next(d.r)
}

func (d *Deque[T]) PopFront() T {
	if d.isEmpty() {
		panic("Deque PopFront deque is empty")
	}
	v := d.buf[d.l]
	d.l = d.next(d.l)
	return v
}

func (d *Deque[T]) PopBack() T {
	if d.isEmpty() {
		panic("Deque PopBack deque is empty")
	}
	d.r = d.prev(d.r)
	return d.buf[d.r]
}

func (d *Deque[T]) At(i int) T {
	if i < 0 || i >= d.size() {
		panic("deque: index out of range")
	}
	return d.buf[(d.l+i)&(len(d.buf)-1)]
}

func (d *Deque[T]) Update(i int, v T) {
	if i < 0 || i >= d.size() {
		panic("deque: index out of range")
	}
	d.buf[(d.l+i)&(len(d.buf)-1)] = v
}

func (d *Deque[T]) Size() int {
	return d.size()
}

func (d *Deque[T]) IsEmpty() bool {
	return d.isEmpty()
}
