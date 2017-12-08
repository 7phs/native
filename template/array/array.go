package array

import "C"
import "bitbucket.org/7phs/fastgotext/wrapper/native"

// template type TArray(A, B, BSize)
type A int
type B C.int

const BSize = C.sizeof_int

type TArray struct {
	native.Array
}

func WithTArray(data []A) *TArray {
	array := NewTArray(uint(len(data)))

	slice := array.Slice()
	for i, v := range data {
		slice[i] = B(v)
	}

	return array
}

func NewTArray(len uint) *TArray {
	return &TArray{
		Array: *native.NewArray(len, BSize),
	}
}

func NewTArrayExt(len uint, pool *native.Pool) *TArray {
	return &TArray{
		Array: *native.NewArrayExt(len, BSize, pool),
	}
}

func (o *TArray) Slice() []B {
	if o.IsEmpty() {
		return []B{}
	}

	return (*[1 << 30]B)(o.Pointer())[:o.Len():o.Len()]
}

func (o *TArray) Marshal() []A {
	if o.IsEmpty() {
		return []A{}
	}

	res := make([]A, 0, o.Len())

	for _, v := range o.Slice() {
		res = append(res, A(v))
	}

	return res
}
