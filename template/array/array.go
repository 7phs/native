package array

import "C"
import "bitbucket.org/7phs/native"

// template type TArray(A, B, BSize)

type A int
type B C.int

const BSize = C.sizeof_int

type TArray struct {
	native.Array
}

type TArrayRec = B

func TArrayPoolKey(dim ...uint) uint {
	return dim[0]
}

func WithTArray(data []A) (array *TArray) {
	array = NewTArray(uint(len(data)))

	slice := array.Slice()
	for i, v := range data {
		slice[i] = B(v)
	}

	return
}

func NewTArray(len uint) *TArray {
	return NewTArrayExt(nil, len)
}

func NewTArrayInterface(pool native.IPool, dim ...uint) interface{} {
	return NewTArrayExt(pool, dim[0])
}

func NewTArrayExt(pool native.IPool, len uint) *TArray {
	return &TArray{Array: *native.NewArrayExt(pool, BSize, len)}
}

func (o *TArray) Slice() []B {
	if o.IsEmpty() {
		return []B{}
	}

	len := o.Dim()[0]

	return (*[1 << 30]B)(o.Pointer())[:len:len]
}

func (o *TArray) Marshal() []A {
	if o.IsEmpty() {
		return []A{}
	}

	res := make([]A, 0, o.Dim()[0])

	for _, v := range o.Slice() {
		res = append(res, A(v))
	}

	return res
}

func (o *TArray) Free() {
	if o.HasPool() {
		o.Pool().Put(o)
	} else {
		o.FreeData()
	}
}
