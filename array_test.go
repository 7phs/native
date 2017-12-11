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

func sliceTestArray(array *Array) []int8 {
	sz := array.Size() * array.itemSize

	return (*[1 << 30]int8)(array.Pointer())[:sz:sz]
}

func generateTestDataArray(slice []int8) []int8 {
	for i := range slice {
		slice[i] = int8(int32(i+3) % 256)
	}

	return slice
}

func sumTestArray(slice []int8) int32 {
	var (
		sum int32 = 0
	)
	for _, v := range slice {
		sum += int32(v)
	}
	return sum
}

func TestArray_Clear(t *testing.T) {
	var (
		expectedSum = []int32{3, 7, 18, 168, 738, 2208, -256}
		len         = uint(1)
	)

	for i, itemSize := range []uint{1, 2, 4, 16, 36, 64, 512} {
		array := NewArray(itemSize, len)

		slice := sliceTestArray(array)
		sum := sumTestArray(generateTestDataArray(slice))

		if sum != expectedSum[i] {
			t.Error("failed to set some values. Sum is", sum, ", but expected is", expectedSum[i])
		}

		array.ClearData()

		if sum = sumTestArray(slice); sum != 0 {
			t.Error("failed to clear data. Sum is", sum)
		}

		if sum = sumTestArray(sliceTestArray(array)); sum != 0 {
			t.Error("failed to clear data (with recreated slice). Sum is", sum)
		}

		array.Free()
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
