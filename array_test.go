package native

import (
	"reflect"
	"testing"
)

func TestNewArray(t *testing.T) {
	array := NewArray(4, 16)
	defer array.Free()
}

func TestNewArrayExt(t *testing.T) {
	var (
		dim           = []uint{16}
		itemSize uint = 4
	)

	pool := NewPool(func(pool IPool) interface{} {
		return NewArrayExt(pool, itemSize, dim...)
	})
	defer pool.FreeData()

	array := NewArrayExt(pool, itemSize, dim...)
	array.Free()
}

func TestArray_Data(t *testing.T) {
	var (
		dim           = []uint{16, 37}
		itemSize uint = 4
	)
	array := NewArray(itemSize)
	defer array.Free()

	if array.Pointer() == nil {
		t.Error("failed to allocate data")
	}

	if array.IsEmpty() {
		t.Error("failed to check allocated data")
	}

	if existDim := array.Dim(); reflect.DeepEqual(existDim, dim) {
		t.Error("failed to get dimention of allocated data. Got\n", existDim, "\n but expected is\n", dim)
	}

	if existItemSize := array.ItemSize(); existItemSize != itemSize {
		t.Error("failed to get item size of allocated data. Got ", existItemSize, ", but expected is ", itemSize)
	}
}

func sliceTestArray(array *Array) []int32 {
	return (*[1 << 30]int32)(array.Pointer())[:array.Size():array.Size()]
}

func generateTestDataArray(slice []int32) []int32 {
	for i := range slice {
		slice[i] = int32(i + 3)
	}

	return slice
}

func sumTestArray(slice []int32) int32 {
	var (
		sum int32 = 0
	)
	for _, v := range slice {
		sum += v
	}
	return sum
}

func TestArray_Clear(t *testing.T) {
	var (
		len      uint = 16
		itemSize uint = 4
	)

	array := NewArray(itemSize, len)
	defer array.Free()

	slice := sliceTestArray(array)
	sum := sumTestArray(generateTestDataArray(slice))
	var expectedSum int32 = 168

	if sum != expectedSum {
		t.Error("failed to set some values. Sum is", sum, ", but expected is", expectedSum)
	}

	array.ClearData()

	if sum = sumTestArray(slice); sum != 0 {
		t.Error("failed to clear data. Sum is", sum)
	}

	if sum = sumTestArray(sliceTestArray(array)); sum != 0 {
		t.Error("failed to clear data (with recreated slice). Sum is", sum)
	}
}

func TestArray_Free(t *testing.T) {
	var (
		dim           = []uint{16, 37, 2}
		itemSize uint = 4
	)

	array := NewArray(itemSize, dim...)
	array.Free()
	array.Free()

	if array.Pointer() != nil {
		t.Error("failed to free data")
	}

	if !array.IsEmpty() {
		t.Error("failed to check empty data")
	}

	if array.Size() != 0 {
		t.Error("failed to clear length while freed data")
	}

	if array.ItemSize() != 0 {
		t.Error("failed to clear item size while freed data")
	}

	// Pooling
	pool := NewPool(func(pool IPool) interface{} {
		return NewArrayExt(pool, itemSize, dim...)
	})

	array = NewArrayExt(pool, itemSize, dim...)
	array.Free()

	if array.Pointer() == nil {
		t.Error("data shouldn't free for pooling")
	}

	if array.IsEmpty() {
		t.Error("data shouldn't empty for pooling")
	}

	if array.Size() == 0 {
		t.Error("data size shouldn't clear for pooling")
	}

	if array.ItemSize() == 0 {
		t.Error("data item size shouldn't clear for pooling")
	}
}

func TestArray_WithPool(t *testing.T) {
	var (
		dim           = []uint{16}
		itemSize uint = 4
	)

	pool := NewPool(func(pool IPool) interface{} {
		return NewArrayExt(pool, itemSize, dim...)
	})
	defer pool.FreeData()

	array := pool.Get().(*Array)

	if sum := sumTestArray(generateTestDataArray(sliceTestArray(array))); sum == 0 {
		t.Error("failed to generate data. Got ", sum)
	}

	array.Free()

	if array.IsEmpty() {
		t.Error("data isn't needing to free for arraypool version")
	}

	if sum := sumTestArray(sliceTestArray(array)); sum != 0 {
		t.Error("data is empty, but should empty")
	}
}
