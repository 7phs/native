package native

import (
    "testing"
)

func TestNewMatrix(t *testing.T) {
	NewMatrix(16, 64, 8)
}

func TestNewMatrixExt(t *testing.T) {
	var (
		rowLen      uint = 16
        colLen uint = 128
        len uint = rowLen * colLen
        		itemSize uint = 4
	)

	pool := NewPool(len, func(itemSize uint, pool *Pool) interface{} {
		return NewMatrixExt(rowLen, colLen, itemSize, pool)
	})

	NewMatrixExt(rowLen, colLen, itemSize, pool)
}
