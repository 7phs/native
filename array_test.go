package native

import (
	"testing"
)

func TestNewArray(t *testing.T) {
	NewArray(16, 4)
}

func TestNewArrayExt(t *testing.T) {
	var (
		len      uint = 16
		itemSize uint = 4
	)

	pool := NewPool(len, func(itemSize uint, pool *Pool) interface{} {
		return NewArrayExt(len, itemSize, pool)
	})

	NewArrayExt(len, itemSize, pool)
}

func TestArrayData(t *testing.T) {
	var (
		len      uint = 16
		itemSize uint = 4
	)
	array := NewArray(len, itemSize)
	if array.Pointer() == nil {
		t.Error("failed to allocate data")
	}

	if array.IsEmpty() {
		t.Error("failed to check allocated data")
	}

	if existLen := array.Len(); existLen != len {
		t.Error("failed to get length of allocated data. Got ", existLen, ", but expected is ", len)
	}

	if existItemSize := array.ItemSize(); existItemSize != itemSize {
		t.Error("failed to get item size of allocated data. Got ", existItemSize, ", but expected is ", itemSize)
	}
}

func TestArrayFree(t *testing.T) {
	var (
		len      uint = 16
		itemSize uint = 4
	)

	array := NewArray(len, itemSize)
	array.Free()
	array.Free()

	if array.Pointer() != nil {
		t.Error("failed to free data")
	}

	if !array.IsEmpty() {
		t.Error("failed to check empty data")
	}

	if array.Len() != 0 {
		t.Error("failed to clear length while freed data")
	}

	if array.ItemSize() != 0 {
		t.Error("failed to clear item size while freed data")
	}

	// Pooling
	pool := NewPool(len, func(itemSize uint, pool *Pool) interface{} {
		return NewArrayExt(len, itemSize, pool)
	})
	array = NewArrayExt(16, 4, pool)

}
