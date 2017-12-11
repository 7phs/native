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
type PoolData interface {
	ItemSize() uint
	Dim() []uint
	ClearData() PoolData
	FreeData()
}

/*
 */
type IPool interface {
	Get() PoolData
	Put(PoolData)
	FreeData()
}

/*
 */
type Status struct {
	status int32
}

func (o *Status) IsFinish() bool {
	return atomic.LoadInt32(&o.status) == POOL_STATUS_FINISH
}

func (o *Status) setFinish() bool {
	return atomic.CompareAndSwapInt32(&o.status, POOL_STATUS_NORMAL, POOL_STATUS_FINISH)
}

/*
 */
func PoolManagerDefaultKey(dim ...uint) (key uint) {
	key = 1
	for _, d := range dim {
		key *= d
	}
	return
}

/*
 */
type PoolManager struct {
	sync.Map
	Status

	itemSize uint
	key      func(...uint) uint
	new      func(IPool, ...uint) interface{}
}

func NewPoolManager(itemSize uint, new func(IPool, ...uint) interface{}) *PoolManager {
	return &PoolManager{
		itemSize: itemSize,
		key:      PoolManagerDefaultKey,
		new:      new,
	}
}

func (o *PoolManager) SetKey(key func(...uint) uint) *PoolManager {
	o.key = key

	return o
}

func (o *PoolManager) getPool(dim ...uint) IPool {
	key := o.key(dim...)

	if pool, ok := o.Load(key); ok {
		return pool.(IPool)
	}

	pool, _ := o.LoadOrStore(key, NewPool(func(pool IPool) interface{} {
		return o.new(pool, dim...)
	}))

	return pool.(IPool)
}

func (o *PoolManager) Get(dim ...uint) interface{} {
	return o.getPool(dim...).Get()
}

func (o *PoolManager) Put(data PoolData) (err error) {
	if o.IsFinish() {
		return
	}

	if data.ItemSize() != o.itemSize {
		err = os.ErrInvalid
		return
	}

	o.getPool(data.Dim()...).Put(data)

	return
}

func (o *PoolManager) FreeData() {
	if !o.setFinish() {
		return
	}

	o.Range(func(_, value interface{}) bool {
		value.(IPool).FreeData()

		return true
	})
}

/*
 */
type Pool struct {
	sync.Pool
	Status

	status int32
}

func NewPool(new func(pool IPool) interface{}) (pool *Pool) {
	pool = &Pool{
		status: POOL_STATUS_NORMAL,
	}
	pool.New = func() interface{} {
		if pool.IsFinish() {
			return nil
		}

		return new(pool)
	}
	return
}

func (o *Pool) Get() PoolData {
	if o.IsFinish() {
		return nil
	}

	return o.Pool.Get().(PoolData)
}

func (o *Pool) Put(data PoolData) {
	if o.IsFinish() {
		return
	}

	data.ClearData()
	o.Pool.Put(data)
}

func (o *Pool) FreeData() {
	if !o.setFinish() {
		return
	}

	for {
		if data := o.Pool.Get(); data == nil {
			break
		} else {
			data.(PoolData).FreeData()
		}
	}
}
