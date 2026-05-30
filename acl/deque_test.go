package acl

import "testing"

func TestDeque_PushBack_GetBack(t *testing.T) {
	d := NewDeque[int]()
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	if got := d.GetBack(); got != 3 {
		t.Errorf("GetBack() = %d, want 3", got)
	}
	if got := d.GetFront(); got != 1 {
		t.Errorf("GetFront() = %d, want 1", got)
	}
	if got := d.Size(); got != 3 {
		t.Errorf("Size() = %d, want 3", got)
	}
}

func TestDeque_PushFront_GetFront(t *testing.T) {
	d := NewDeque[int]()
	d.PushFront(1)
	d.PushFront(2)
	d.PushFront(3)

	if got := d.GetFront(); got != 3 {
		t.Errorf("GetFront() = %d, want 3", got)
	}
	if got := d.GetBack(); got != 1 {
		t.Errorf("GetBack() = %d, want 1", got)
	}
}

func TestDeque_PopFrontPopBack(t *testing.T) {
	d := NewDeque[int]()
	for i := 1; i <= 5; i++ {
		d.PushBack(i)
	}

	if got := d.PopFront(); got != 1 {
		t.Errorf("PopFront() = %d, want 1", got)
	}
	if got := d.PopBack(); got != 5 {
		t.Errorf("PopBack() = %d, want 5", got)
	}
	if got := d.GetFront(); got != 2 {
		t.Errorf("GetFront() = %d, want 2", got)
	}
	if got := d.GetBack(); got != 4 {
		t.Errorf("GetBack() = %d, want 4", got)
	}
	if got := d.Size(); got != 3 {
		t.Errorf("Size() = %d, want 3", got)
	}
}

func TestDeque_AtUpdate(t *testing.T) {
	d := NewDeque[int]()
	for _, v := range []int{10, 20, 30, 40} {
		d.PushBack(v)
	}

	for i, want := range []int{10, 20, 30, 40} {
		if got := d.At(i); got != want {
			t.Errorf("At(%d) = %d, want %d", i, got, want)
		}
	}

	d.Update(2, 99)
	if got := d.At(2); got != 99 {
		t.Errorf("At(2) after Update = %d, want 99", got)
	}
}

func TestDeque_MixedPushPop(t *testing.T) {
	d := NewDeque[int]()
	d.PushBack(1)
	d.PushFront(2)
	d.PushBack(3)
	d.PushFront(4) // [4, 2, 1, 3]

	if got := d.GetFront(); got != 4 {
		t.Errorf("GetFront() = %d, want 4", got)
	}
	if got := d.GetBack(); got != 3 {
		t.Errorf("GetBack() = %d, want 3", got)
	}

	want := []int{4, 2, 1, 3}
	for i, w := range want {
		if got := d.At(i); got != w {
			t.Errorf("At(%d) = %d, want %d", i, got, w)
		}
	}
}

func TestDeque_Grow(t *testing.T) {
	d := NewDeque[int]()
	n := DEQUE_CAP + 100
	for i := range n {
		d.PushBack(i)
	}
	if got := d.Size(); got != n {
		t.Errorf("Size() = %d, want %d", got, n)
	}
	if got := d.GetFront(); got != 0 {
		t.Errorf("GetFront() = %d, want 0", got)
	}
	if got := d.GetBack(); got != n-1 {
		t.Errorf("GetBack() = %d, want %d", got, n-1)
	}
	for i := range n {
		if got := d.PopFront(); got != i {
			t.Errorf("PopFront[%d] = %d, want %d", i, got, i)
		}
	}
	if !d.IsEmpty() {
		t.Error("deque should be empty after popping all elements")
	}
}

func TestDeque_PanicOnEmptyPopFront(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on PopFront from empty deque")
		}
	}()
	d := NewDeque[int]()
	d.PopFront()
}

func TestDeque_PanicOnEmptyPopBack(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on PopBack from empty deque")
		}
	}()
	d := NewDeque[int]()
	d.PopBack()
}
