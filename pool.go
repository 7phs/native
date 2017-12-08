package native

import (
	"sync"
)

/*
 */
type PoolId struct {
	Kind     uint
	ItemSize uint
	New      func(*Pool) interface{}
}

/*
 */
type PoolManager struct {
	sync.Map
}

func NewPoolManager() *PoolManager {
	return &PoolManager{}
}

func (o *PoolManager) getPool(key *PoolId) *Pool {
	if pool, ok := o.Load(key); ok {
		return pool.(*Pool)
	}

	pool, _ := o.LoadOrStore(key, NewPool(key.New))

	return pool.(*Pool)
}

func (o *PoolManager) Get(key *PoolId) interface{} {
	return o.getPool(key).Get()
}

func (o *PoolManager) Put(key *PoolId, data PoolData) {
	o.getPool(key).Put(data)
}

/*
 */
type PoolData interface {
	ClearData() PoolData
	FreeData()
}

/*
 */
type Pool struct {
	sync.Pool
}

func NewPool(new func(*Pool) interface{}) (pool *Pool) {
	pool = &Pool{}
	pool.New = func() interface{} {
		return new(pool)
	}
	return
}

func (o *Pool) Get() PoolData {
	return o.Pool.Get().(PoolData).ClearData()
}

func (o *Pool) Put(data PoolData) {
	o.Pool.Put(data)
}
