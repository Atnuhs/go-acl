package acl

import (
	"sort"
	"testing"
)

func TestPQ(t *testing.T) {
	p := NewPriorityQueue[int]()

	values := []int{3, 4, 1, 2, 5, 8, 7, 6, -121}
	for _, v := range values {
		p.Push(v)
	}
	got := make([]int, p.Size())
	for i := range got {
		got[i] = p.Pop()
	}
	sort.Ints(values)
	for i := range got {
		if got[i] != values[i] {
			t.Errorf("invalid at index %d want %d, got %d", i, values[i], got[i])
		}
	}
}

func TestPQ_IsEmpty(t *testing.T) {
	p := NewPriorityQueue[int]()
	if !p.IsEmpty() {
		t.Error("new queue should be empty")
	}
	p.Push(1)
	if p.IsEmpty() {
		t.Error("queue with element should not be empty")
	}
	p.Pop()
	if !p.IsEmpty() {
		t.Error("queue should be empty after pop")
	}
}

func TestPQ_Size(t *testing.T) {
	p := NewPriorityQueue[int]()
	for i := range 5 {
		if p.Size() != i {
			t.Errorf("Size() = %d, want %d", p.Size(), i)
		}
		p.Push(i)
	}
}

func TestPQ_Peek(t *testing.T) {
	p := NewPriorityQueue[int]()
	p.Push(5)
	p.Push(3)
	p.Push(7)

	if got := p.Peek(); got != 3 {
		t.Errorf("Peek() = %d, want 3", got)
	}
	if p.Size() != 3 {
		t.Error("Peek should not remove element")
	}
}

func TestPQ_MaxHeap(t *testing.T) {
	p := NewPriorityQueueWithLessFunc(func(a, b int) bool { return a > b })

	values := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, v := range values {
		p.Push(v)
	}

	expected := make([]int, len(values))
	copy(expected, values)
	sort.Sort(sort.Reverse(sort.IntSlice(expected)))

	for i, want := range expected {
		got := p.Pop()
		if got != want {
			t.Errorf("Pop[%d] = %d, want %d", i, got, want)
		}
	}
}

func TestPQ_PanicOnEmptyPop(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Pop from empty queue")
		}
	}()
	p := NewPriorityQueue[int]()
	p.Pop()
}

func TestPQ_PanicOnEmptyPeek(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Peek from empty queue")
		}
	}()
	p := NewPriorityQueue[int]()
	p.Peek()
}
