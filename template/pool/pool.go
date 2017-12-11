package pool

import "C"
import "bitbucket.org/7phs/native"

// template type TPoolManager(A, BSize, KeyDimention, ItemNew)

type A struct{ native.Array }

const BSize = C.sizeof_int

var KeyDimention = native.PoolManagerDefaultKey

var ItemNew = func(pool native.IPool, dim ...uint) interface{} {
	return &A{Array: *native.NewArrayExt(pool, C.sizeof_int, dim...)}
}

type TPoolManager struct {
	native.PoolManager
}

func NewTPoolManager() *TPoolManager {
	return &TPoolManager{
		PoolManager: *native.NewPoolManager(BSize, ItemNew).SetKey(KeyDimention),
	}
}

func (o *TPoolManager) Get(dim ...uint) *A {
	return o.PoolManager.Get(dim...).(*A)
}
