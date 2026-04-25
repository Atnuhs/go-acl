package acl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestAns(t *testing.T) {
	testOut := new(bytes.Buffer)
	Out = bufio.NewWriter(testOut)
	testCases := map[string]struct {
		data     []interface{}
		expected string
	}{
		"only int":    {data: []interface{}{1, 2, 3}, expected: "1 2 3\n"},
		"only string": {data: []interface{}{"a", "b", "c"}, expected: "a b c\n"},
		"only []int":  {data: []interface{}{[]int{1, 2, 3}}, expected: "1 2 3\n"},
		"combined":    {data: []interface{}{1, 2, 3, "4", "a", []int{5, 6, 7}}, expected: "1 2 3 4 a 5 6 7\n"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testOut.Reset()
			Ans(tc.data...)
			Out.Flush()
			actual := testOut.String()
			if tc.expected != actual {
				t.Errorf("expected: %q, but got: %q", tc.expected, actual)
			}
		})
	}
}

func TestYesNo(t *testing.T) {
	testOut := new(bytes.Buffer)
	Out = bufio.NewWriter(testOut)

	tests := map[string]struct {
		input bool
		want  string
	}{
		"true":  {input: true, want: "Yes\n"},
		"false": {input: false, want: "No\n"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testOut.Reset()
			YesNo(tc.input)
			Out.Flush()
			actual := testOut.String()
			if tc.want != actual {
				t.Errorf("expected: %q, but got: %q", tc.want, actual)
			}
		})
	}
}

func TestYesNoFunc(t *testing.T) {
	testOut := new(bytes.Buffer)
	Out = bufio.NewWriter(testOut)

	tests := map[string]struct {
		inputFunc func() bool
		want      string
	}{
		"function returns true": {
			inputFunc: func() bool { return true },
			want:      "Yes\n",
		},
		"function returns false": {
			inputFunc: func() bool { return false },
			want:      "No\n",
		},
		"complex function true": {
			inputFunc: func() bool { return 5 > 3 },
			want:      "Yes\n",
		},
		"complex function false": {
			inputFunc: func() bool { return 2 > 5 },
			want:      "No\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testOut.Reset()
			YesNoFunc(tc.inputFunc)
			Out.Flush()
			actual := testOut.String()
			if tc.want != actual {
				t.Errorf("expected: %q, but got: %q", tc.want, actual)
			}
		})
	}
}

func BenchmarkOutputToOut(b *testing.B) {
	text := strings.Repeat("a", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintln(Out, text)
	}
}

func BenchmarkOutputToDiscard(b *testing.B) {
	text := strings.Repeat("a", 100)
	Discard := bufio.NewWriter(io.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintln(Discard, text)
	}
}

func TestKeyInts(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	want := "1 2 3 4 5"
	got := KeyInts(arr)

	if got != Key(want) {
		t.Errorf("want %s but got %s", want, got)
	}
}

func TestKey_ToInts(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	key := KeyInts(arr)
	got := key.ToInts()
	if len(got) != len(arr) {
		t.Fatalf("length not match want %d but got %d", len(arr), len(got))
	}
	for i, wantV := range arr {
		gotV := got[i]
		if wantV != gotV {
			t.Errorf("idx %d: want %d got %d", i, wantV, gotV)
		}
	}
}

func TestSort(t *testing.T) {
	arr := []int{2, 3, 4, 5, 1}
	want := []int{1, 2, 3, 4, 5}

	got := Sort(arr)
	if len(got) != len(want) {
		t.Fatalf("length not match want %d but got %d", len(want), len(got))
	}
	for i, wantV := range want {
		gotV := got[i]
		if wantV != gotV {
			t.Errorf("idx %d: want %d got %d", i, wantV, gotV)
		}
	}

	idxs := SortIdx(arr)
	if len(idxs) != len(want) {
		t.Fatalf("length not match want %d but got %d", len(want), len(idxs))
	}
	for i, wantV := range want {
		gotV := arr[idxs[i]]
		if wantV != gotV {
			t.Errorf("idx %d: want %d got %d", i, wantV, gotV)
		}
	}

	type st struct{ v int }
	arr2 := make([]st, len(arr))
	less := LessFunc[st](func(a, b st) bool { return a.v < b.v })
	for i := range arr2 {
		arr2[i] = st{v: arr[i]}
	}
	got2 := SortF(arr2, less)

	if len(got2) != len(want) {
		t.Fatalf("length not match want %d but got %d", len(want), len(got2))
	}
	for i, wantV := range want {
		gotV := got2[i].v
		if wantV != gotV {
			t.Errorf("idx %d: want %d got %d", i, wantV, gotV)
		}
	}

	idxs2 := SortIdxF(arr2, less)
	if len(idxs2) != len(want) {
		t.Fatalf("length not match want %d but got %d", len(want), len(idxs))
	}
	for i, wantV := range want {
		gotV := arr2[idxs2[i]].v
		if wantV != gotV {
			t.Errorf("idx %d: want %d got %d", i, wantV, gotV)
		}
	}
}
