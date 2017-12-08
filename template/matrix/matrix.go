package matrix

import "C"
import "bitbucket.org/7phs/fastgotext/wrapper/native"

// template type TMatrix(AM, BM, BMSize)
type AM int
type BM C.int

const BMSize = C.sizeof_int

type TMatrix struct {
	native.Matrix
}

func NewTMatrix(rowLen, colLen uint) *TMatrix {
	return &TMatrix{
		Matrix: *native.NewMatrix(rowLen, colLen, BMSize),
	}
}

func NewTMatrixExt(rowLen, colLen uint, pool *native.Pool) *TMatrix {
	return &TMatrix{
		Matrix: *native.NewMatrixExt(rowLen, colLen, BMSize, pool),
	}
}

func (o *TMatrix) Slice() [][]BM {
	if o.IsEmpty() {
		return [][]BM{}
	}

	var (
		rowLen, colLen = o.Len()
		sz             = rowLen * colLen
		col            = int(colLen)
		start          int

		data   = (*[1 << 30]BM)(o.Pointer())[:sz:sz]
		slices = make([][]BM, rowLen)
	)

	for i := range slices {
		start = col * i
		slices[i] = data[start : start+col]
	}

	return slices
}

func (o *TMatrix) Marshal() [][]AM {
	if o.IsEmpty() {
		return [][]AM{}
	}
	rowLen, _ := o.Len()

	slices := make([][]AM, 0, rowLen)

	for _, rowN := range o.Slice() {
		row := make([]AM, 0, rowLen)

		for _, v := range rowN {
			row = append(row, AM(v))
		}

		slices = append(slices, row)
	}

	return slices
}
