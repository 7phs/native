package native

// #include <stdlib.h>
// #include <string.h>
//
// void* native_malloc(size_t len, size_t item_size) {
//     return(malloc(len * item_size));
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
import (
	"unsafe"
)

type clearSlice struct {
	clearType int
	i64       []int64
	i32       []int32
	i16       []int16
	i8        []int8
}

func (o *clearSlice) init(data *C.void, clearLen uint, clearType int) {
	switch o.clearType = clearType; o.clearType {
	case 8:
		o.i64 = (*[1 << 30]int64)(unsafe.Pointer(data))[:clearLen:clearLen]
	case 4:
		o.i32 = (*[1 << 30]int32)(unsafe.Pointer(data))[:clearLen:clearLen]
	case 2:
		o.i16 = (*[1 << 30]int16)(unsafe.Pointer(data))[:clearLen:clearLen]
	default:
		o.i8 = (*[1 << 30]int8)(unsafe.Pointer(data))[:clearLen:clearLen]
	}
}

func (o *clearSlice) Clear() {
	switch o.clearType {
	case 8:
		for i := range o.i64 {
			o.i64[i] = 0
		}
	case 4:
		for i := range o.i32 {
			o.i32[i] = 0
		}
	case 2:
		for i := range o.i16 {
			o.i16[i] = 0
		}
	default:
		for i := range o.i8 {
			o.i8[i] = 0
		}
	}
}

type Array struct {
	data       *C.void
	dim        []uint
	size       uint
	itemSize   uint
	clearSlice clearSlice

	pool IPool
}

func NewArray(itemSize uint, dim ...uint) *Array {
	return NewArrayExt(nil, itemSize, dim...)
}

func newArrayExtParams(itemSize uint, dim ...uint) (size uint, clearLen uint, clearType int) {
	size = 1
	for _, d := range dim {
		size *= d
	}

	byteSize := size * itemSize
	if byteSize%8 == 0 {
		clearType = 8
		clearLen = byteSize / 8
	} else if byteSize%4 == 0 {
		clearType = 4
		clearLen = byteSize / 4
	} else if byteSize%2 == 0 {
		clearType = 2
		clearLen = byteSize / 2
	}

	return
}

func NewArrayExt(pool IPool, itemSize uint, dim ...uint) *Array {
	size, clearLen, clearType := newArrayExtParams(itemSize, dim...)

	return (&Array{
		data:     (*C.void)(C.native_malloc(C.size_t(size), C.size_t(itemSize))),
		dim:      dim,
		size:     size,
		itemSize: itemSize,
		pool:     pool,
	}).initClearSlice(clearLen, clearType)
}

func (o *Array) initClearSlice(clearLen uint, clearType int) *Array {
	o.clearSlice.init(o.data, clearLen, clearType)

	return o
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
	o.clearSlice.Clear()

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
