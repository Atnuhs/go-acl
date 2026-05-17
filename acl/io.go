package acl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	BufferSize = 1 << 20
)

var (
	Out = bufio.NewWriterSize(os.Stdout, BufferSize)
	Dbg = bufio.NewWriterSize(os.Stderr, BufferSize)

	in = bufio.NewReaderSize(os.Stdin, BufferSize)
)

// --- Input ---

// readToken は空白をスキップして次のトークンのバイト列を返す
func readToken() []byte {
	c, err := in.ReadByte()
	for err == nil && c <= ' ' {
		c, err = in.ReadByte()
	}
	if err != nil {
		return nil
	}
	buf := make([]byte, 0, 16)
	for err == nil && c > ' ' {
		buf = append(buf, c)
		c, err = in.ReadByte()
	}
	return buf
}

// S は文字列を読み込む
func S() string {
	return string(readToken())
}

// B は文字列を[]byteとして読み込む
func B() []byte {
	return readToken()
}

func F() float64 {
	ret, _ := strconv.ParseFloat(string(readToken()), 64)
	return ret
}

// I は整数を読み込む
func I() int {
	c, err := in.ReadByte()
	for err == nil && c <= ' ' {
		c, err = in.ReadByte()
	}
	if err != nil {
		return 0
	}
	neg := c == '-'
	if neg {
		c, err = in.ReadByte()
	}
	n := 0
	for err == nil && '0' <= c && c <= '9' {
		n = n*10 + int(c-'0')
		c, err = in.ReadByte()
	}
	if neg {
		return -n
	}
	return n
}

// II は整数を2つ読み込む
func II() (int, int) {
	return I(), I()
}

// III は整数を3つ読み込む
func III() (int, int, int) {
	return I(), I(), I()
}

// IIII は整数を4つ読み込む
func IIII() (int, int, int, int) {
	return I(), I(), I(), I()
}

// Is は整数をn個読み込む
func Is(n int) []int {
	ret := make([]int, n)
	for i := range ret {
		ret[i] = I()
	}
	return ret
}

// IIs はi行目がa_i, b_iとなっている入力を受け取り、配列a, bを返す
func IIs(n int) ([]int, []int) {
	a := make([]int, n)
	b := make([]int, n)
	for i := range a {
		a[i], b[i] = II()
	}
	return a, b
}

// IIIs はi行目がa_i, b_i, c_iとなっている入力を受け取り、配列a, b, cを返す
func IIIs(n int) ([]int, []int, []int) {
	a := make([]int, n)
	b := make([]int, n)
	c := make([]int, n)
	for i := range a {
		a[i], b[i], c[i] = III()
	}
	return a, b, c
}

// IIIIs はi行目がa_i, b_i, c_i, d_iとなっている入力を受け取り、配列a, b, c, dを返す
func IIIIs(n int) ([]int, []int, []int, []int) {
	a := make([]int, n)
	b := make([]int, n)
	c := make([]int, n)
	d := make([]int, n)
	for i := range a {
		a[i], b[i], c[i], d[i] = IIII()
	}
	return a, b, c, d
}

// ReadGrid は整数をh行w列の配列として読み込む
func ReadGrid(h, w int) [][]int {
	ret := L2[int](h, w)
	for i := range ret {
		for j := range ret[i] {
			ret[i][j] = I()
		}
	}
	return ret
}

// Ss は文字列をn個読み込む
func Ss(n int) []string {
	ret := L1[string](n)
	for i := range ret {
		ret[i] = S()
	}
	return ret
}

// Bs は文字列をn個読み込む
func Bs(n int) [][]byte {
	ret := make([][]byte, n)
	for i := range ret {
		ret[i] = B()
	}
	return ret
}

func Fs(n int) []float64 {
	ret := make([]float64, n)
	for i := range ret {
		ret[i] = F()
	}
	return ret
}

// --- Output ---

func writeInt(v int) {
	_, _ = Out.WriteString(strconv.Itoa(v))
}

func writeString(v string) {
	_, _ = Out.WriteString(v)
}

func writeFloat64(v float64) {
	_, _ = Out.WriteString(strconv.FormatFloat(v, 'f', 14, 64))
}

func writeByte(v byte) {
	_ = Out.WriteByte(v)
}

func writeByteSlice(v []byte) {
	_, _ = Out.Write(v)
}

func writeByteGrid(g [][]byte) {
	for i, row := range g {
		if i > 0 {
			_ = Out.WriteByte('\n')
		}
		writeByteSlice(row)
	}
}

func writeSlice[T any](s []T, writeElem func(T)) {
	for i, x := range s {
		if i > 0 {
			_ = Out.WriteByte(' ')
		}
		writeElem(x)
	}
}

func writeGrid[T any](g [][]T, writeElem func(T)) {
	for i, row := range g {
		if i > 0 {
			_ = Out.WriteByte('\n')
		}
		writeSlice(row, writeElem)
	}
}

func writeOne(v any) {
	switch x := v.(type) {
	case int:
		writeInt(x)
	case string:
		writeString(x)
	case float64:
		writeFloat64(x)
	case byte:
		writeByte(x)
	case []int:
		writeSlice(x, writeInt)
	case []string:
		writeSlice(x, writeString)
	case []float64:
		writeSlice(x, writeFloat64)
	case []byte:
		writeByteSlice(x)
	case [][]int:
		writeGrid(x, writeInt)
	case [][]string:
		writeGrid(x, writeString)
	case [][]float64:
		writeGrid(x, writeFloat64)
	case [][]byte:
		writeByteGrid(x)
	default:
		fmt.Fprint(Out, x)
	}
}

// Ans は出力を行う
func Ans(args ...any) {
	for i, arg := range args {
		if i > 0 {
			_ = Out.WriteByte(' ')
		}
		writeOne(arg)
	}
	_ = Out.WriteByte('\n')
}

// Yes は"Yes"を出力する
func Yes() {
	Ans("Yes")
}

// No は"No"を出力する
func No() {
	Ans("No")
}

// YesNo は条件に応じてYesまたはNoを出力する
func YesNo(b bool) {
	if b {
		Yes()
	} else {
		No()
	}
}

// YesNoFunc は関数の結果に応じてYesまたはNoを出力する
func YesNoFunc(f func() bool) {
	YesNo(f())
}
