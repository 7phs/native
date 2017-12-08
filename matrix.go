package native

type Matrix struct {
	Array

	rowLen uint
	colLen uint
}

func NewMatrix(rowLen, colLen, itemSize uint) *Matrix {
	return NewMatrixExt(rowLen, colLen, itemSize, nil)
}

func NewMatrixExt(rowLen, colLen, itemSize uint, pool *Pool) *Matrix {
	return &Matrix{
		Array: *NewArrayExt(rowLen*colLen, itemSize, pool),

		rowLen: rowLen,
		colLen: colLen,
	}
}

func (o *Matrix) Len() (uint, uint) {
	return o.rowLen, o.colLen
}
