package native

// #include <stdlib.h>
// #include <string.h>
//
// void* native_malloc(size_t len, size_t item_size) {
//     return((int*)malloc(len * item_size));
// }
//
// void native_clear(void* data, size_t len, size_t item_size) {
//     memset(data, 0, len * item_size);
// }
//
// void native_free(void* data) {
//     free(data);
// }
import "C"
import "unsafe"

type Array struct {
	data     *C.void
	dim      []uint
	size     uint
	itemSize uint

	pool IPool
}

func NewArray(itemSize uint, dim ...uint) *Array {
	return NewArrayExt(nil, itemSize, dim...)
}

func NewArrayExt(pool IPool, itemSize uint, dim ...uint) *Array {
	size := UintIterator(dim).Mul()

	return &Array{
		data:     (*C.void)(C.native_malloc(C.size_t(size), C.size_t(itemSize))),
		dim:      dim,
		size:     size,
		itemSize: itemSize,
		pool:     pool,
	}
}

func (o *Array) IsEmpty() bool {
	return o.data == nil
}

func (o *Array) Pointer() unsafe.Pointer {
	return unsafe.Pointer(o.data)
}

func (o *Array) Dim() []uint {
	return o.dim
}

func (o *Array) Size() uint {
	return o.size
}

func (o *Array) ItemSize() uint {
	return o.itemSize
}

func (o *Array) Clear() *Array {
	C.native_clear(o.Pointer(), C.size_t(o.size), C.size_t(o.itemSize))

	return o
}

func (o *Array) ClearData() PoolData {
	return o.Clear()
}

func (o *Array) HasPool() bool {
	return o.pool != nil
}

func (o *Array) Pool() IPool {
	return o.pool
}

func (o *Array) Free() {
	if o.HasPool() {
		o.Pool().Put(o)
	} else {
		o.FreeData()
	}
}

func (o *Array) FreeData() {
	if o.IsEmpty() {
		return
	}

	C.native_free(unsafe.Pointer(o.data))

	o.data = nil
	o.dim = []uint{}
	o.size = 0
	o.itemSize = 0
}
