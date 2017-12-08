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
	len      uint
	itemSize uint

	pool *Pool
}

func NewArray(len, itemSize uint) *Array {
	return NewArrayExt(len, itemSize, nil)
}

func NewArrayExt(len, itemSize uint, pool *Pool) *Array {
	return &Array{
		data:     (*C.void)(C.native_malloc(C.size_t(len), C.size_t(itemSize))),
		len:      len,
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

func (o *Array) Len() uint {
	return o.len
}

func (o *Array) ItemSize() uint {
	return o.itemSize
}

func (o *Array) Clear() *Array {
	C.native_clear(o.Pointer(), C.size_t(o.len), C.size_t(o.itemSize))

	return o
}

func (o *Array) ClearData() PoolData {
	return o.Clear()
}

func (o *Array) Free() {
	if o.pool != nil {
		o.pool.Put(o)
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
	o.len = 0
	o.itemSize = 0
}
