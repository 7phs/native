package matrix

import (
	"reflect"
	"testing"

	"bitbucket.org/7phs/native"
)

func TestTMatrixPoolKey(t *testing.T) {
	if key := TMatrixPoolKey(23, 45, 64, 48); key != 23*45 {
		t.Error("failed to get multiplicayion of the first two dimentions in a list. Got ", key, ", but expected is ", 23*45)
	}
}

func TestNewTMatrix(t *testing.T) {
	array := NewTMatrix(16, 37)
	defer array.Free()
}

func TestNewTMatrixExt(t *testing.T) {
	pool := native.NewPool(func(pool native.IPool) interface{} {
		return NewTMatrixExt(pool, 27, 34)
	})
	defer pool.FreeData()

	matrix := pool.Get().(*TMatrix)
	slice := matrix.Slice()
	for i, row := range slice {
		for j := range row {
			slice[i][j] = BM(i) + BM(j)*3
		}
	}

	sum := 0
	for _, row := range matrix.Slice() {
		for _, v := range row {
			sum += int(v)
		}
	}

	if sum == 0 {
		t.Error("failed to set data to array. Sum is empty")
	}

	matrix.Free()
}

func TestNewTMatrixInterface(t *testing.T) {
	pool := native.NewPool(func(pool native.IPool) interface{} {
		return NewTMatrixInterface(pool, 27, 34, 67)
	})
	defer pool.FreeData()

	matrix := pool.Get().(*TMatrix)
	slice := matrix.Slice()
	for i, row := range slice {
		for j := range row {
			slice[i][j] = BM(i) + BM(j)*3
		}
	}

	sum := 0
	for _, row := range matrix.Slice() {
		for _, v := range row {
			sum += int(v)
		}
	}

	if sum == 0 {
		t.Error("failed to set data to array. Sum is empty")
	}

	matrix.Free()
}

func TestWithTMatrix(t *testing.T) {
	expected := [][]AM{
		{0, 2, 4, 6, 8},
		{10, 12, 14, 16, 18},
		{20, 22, 24, 26, 28},
	}

	d := WithTMatrix(expected)
	defer d.Free()

	if exist := d.Marshal(); !reflect.DeepEqual(exist, expected) {
		t.Error("failed to make copy as native matrix. Got\n", exist, "\n but expected is\n", expected)
	}

	WithTMatrix([][]AM{})
}

func TestTMatrix_Clear(t *testing.T) {
	matrix := NewTMatrix(16, 37)

	if len := len(matrix.Slice()); len == 0 {
		t.Error("failed to init tmatrix. Slice length is 0")
	}

	if len := len(matrix.Marshal()); len == 0 {
		t.Error("failed to init tmatrix. Marshaled slice length is 0")
	}

	matrix.Free()

	if len := len(matrix.Slice()); len != 0 {
		t.Error("failed to clear tmatrix. Slice length is ", len)
	}

	if len := len(matrix.Marshal()); len != 0 {
		t.Error("failed to clear tmatrix. Marshaled slice length is ", len)
	}
}
