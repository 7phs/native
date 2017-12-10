package matrix

import "C"
import "bitbucket.org/7phs/native"

// template type TMatrix(AM, BM, BMSize)
type AM int
type BM C.int

const BMSize = C.sizeof_int

type TMatrix struct {
	native.Array
}

type TMatrixRec = BM

func WithTMatrix(data [][]AM) (matrix *TMatrix) {
	var (
		rowLen = uint(len(data))
		colLen = func() uint {
			if rowLen > 0 {
				return uint(len(data[0]))
			} else {
				return 0
			}
		}()
	)

	matrix = NewTMatrix(rowLen, colLen)

	slice := matrix.Slice()
	for i, row := range data {
		for j, v := range row {
			slice[i][j] = BM(v)
		}
	}

	return
}

func NewTMatrix(rowLen, colLen uint) *TMatrix {
	return NewTMatrixExt(nil, rowLen, colLen)
}

func NewTMatrixInterface(pool native.IPool, dim ...uint) interface{} {
	return NewTMatrixExt(pool, dim[0], dim[1])
}

func NewTMatrixExt(pool native.IPool, rowLen, colLen uint) *TMatrix {
	return &TMatrix{Array: *native.NewArrayExt(pool, BMSize, rowLen, colLen)}
}

func (o *TMatrix) Slice() (slices [][]BM) {
	if o.IsEmpty() {
		return [][]BM{}
	}

	var (
		sz     = o.Size()
		rowLen = o.Dim()[0]
		colLen = int(o.Dim()[1])
		start  int

		data = (*[1 << 30]BM)(o.Pointer())[:sz:sz]
	)

	slices = make([][]BM, rowLen)
	for i := range slices {
		start = colLen * i
		slices[i] = data[start : start+colLen]
	}

	return
}

func (o *TMatrix) Marshal() (slices [][]AM) {
	if o.IsEmpty() {
		return [][]AM{}
	}

	rowLen := o.Dim()[0]
	colLen := o.Dim()[1]

	slices = make([][]AM, 0, rowLen)

	for _, rowN := range o.Slice() {
		row := make([]AM, 0, colLen)

		for _, v := range rowN {
			row = append(row, AM(v))
		}

		slices = append(slices, row)
	}

	return
}

func (o *TMatrix) Free() {
	if o.HasPool() {
		o.Pool().Put(o)
	} else {
		o.FreeData()
	}
}
