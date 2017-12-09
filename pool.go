package native

import (
	"os"
	"sync"
	"sync/atomic"
)

const (
	POOL_STATUS_NORMAL int32 = iota
	POOL_STATUS_FINISH
)

/*
 */
type PoolManager struct {
	sync.Map

	itemSize uint
	new      func(uint, *Pool) interface{}
}

func NewPoolManager(itemSize uint, new func(uint, *Pool) interface{}) *PoolManager {
	return &PoolManager{
		itemSize: itemSize,
		new:      new,
	}
}

func (o *PoolManager) getPool(len uint) *Pool {
	if pool, ok := o.Load(len); ok {
		return pool.(*Pool)
	}

	pool, _ := o.LoadOrStore(len, NewPool(len, o.new))

	return pool.(*Pool)
}

func (o *PoolManager) Get(len uint) interface{} {
	return o.getPool(len).Get()
}

func (o *PoolManager) Put(data PoolData) (err error) {
	if data.ItemSize() != o.itemSize {
		err = os.ErrInvalid
		return
	}

	o.getPool(data.Len()).Put(data)

	return
}

/*
 */
type PoolData interface {
	ItemSize() uint
	Len() uint
	ClearData() PoolData
	FreeData()
}

/*
 */
type Pool struct {
	sync.Pool

	status int32
}

func NewPool(len uint, new func(itemSize uint, pool *Pool) interface{}) (pool *Pool) {
	pool = &Pool{
		status: POOL_STATUS_NORMAL,
	}
	pool.New = func() interface{} {
		if pool.isFinish() {
			return nil
		}

		return new(len, pool)
	}
	return
}

func (o *Pool) isFinish() bool {
	return atomic.LoadInt32(&o.status) == POOL_STATUS_FINISH
}

func (o *Pool) setFinish() bool {
	return atomic.CompareAndSwapInt32(&o.status, POOL_STATUS_NORMAL, POOL_STATUS_FINISH)
}

func (o *Pool) Get() PoolData {
	if o.isFinish() {
		return nil
	}

	return o.Pool.Get().(PoolData).ClearData()
}

func (o *Pool) Put(data PoolData) {
	if o.isFinish() {
		return
	}

	o.Pool.Put(data)
}

func (o *Pool) FreeData() {
	if !o.setFinish() {
		return
	}

	for {

		if data := o.Pool.Get().(PoolData); data == nil {
			break
		} else {
			data.FreeData()
		}
	}
}
