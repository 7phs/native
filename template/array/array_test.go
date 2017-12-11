package array

import (
	"bitbucket.org/7phs/native"
	"reflect"
	"testing"
)

func TestTArrayPoolKey(t *testing.T) {
	if key := TArrayPoolKey(23, 45, 64); key != 23 {
		t.Error("failed to get the first dimention in a list. Got ", key, ", but expected is ", 23)
	}
}

func TestNewTArray(t *testing.T) {
	array := NewTArray(16)
	defer array.Free()
}

func TestNewTArrayExt(t *testing.T) {
	pool := native.NewPool(func(pool native.IPool) interface{} {
		return NewTArrayExt(pool, 32)
	})
	defer pool.FreeData()

	array := pool.Get().(*TArray)
	slice := array.Slice()
	for i := range slice {
		slice[i] = B(i) + B(i)*3
	}

	sum := 0
	for _, v := range array.Slice() {
		sum += int(v)
	}

	if sum == 0 {
		t.Error("failed to set data to array. Sum is empty")
	}

	array.Free()
}

func TestNewTArrayInterface(t *testing.T) {
	pool := native.NewPool(func(pool native.IPool) interface{} {
		return NewTArrayInterface(pool, 32)
	})
	defer pool.FreeData()

	array := pool.Get().(*TArray)
	slice := array.Slice()
	for i := range slice {
		slice[i] = B(i) + B(i)*3
	}

	sum := 0
	for _, v := range array.Slice() {
		sum += int(v)
	}

	if sum == 0 {
		t.Error("failed to set data to array. Sum is empty")
	}

	array.Free()
}

func TestWithTArray(t *testing.T) {
	expected := []A{0, 2, 4, 6, 8}

	d := WithTArray(expected)
	defer d.Free()

	if exist := d.Marshal(); !reflect.DeepEqual(exist, expected) {
		t.Error("failed to make copy as native array. Got\n", exist, "\n but expected is\n", expected)
	}
}

func TestTArray_Clear(t *testing.T) {
	array := NewTArray(16)

	if len := len(array.Slice()); len == 0 {
		t.Error("failed to init tarray. Slice length is 0")
	}

	if len := len(array.Marshal()); len == 0 {
		t.Error("failed to init tarray. Marshaled slice length is 0")
	}

	array.Free()

	if len := len(array.Slice()); len != 0 {
		t.Error("failed to clear tarray. Slice length is ", len)
	}

	if len := len(array.Marshal()); len != 0 {
		t.Error("failed to clear tarray. Marshaled slice length is ", len)
	}
}
