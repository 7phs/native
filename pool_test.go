package native

import (
	"fmt"
	"strings"
	"testing"
)

type testNativeRec struct {
	data     string
	itemSize uint
	dim      []uint

	free bool

	pool IPool
}

func newTestNativeRec(data string, itemSize uint, dim ... uint) *testNativeRec {
	return newTestNativeRecExt(nil, data, itemSize, dim...)
}

func newTestNativeRecExt(pool IPool, data string, itemSize uint, dim ... uint) *testNativeRec {
	return &testNativeRec{
		data:     data,
		itemSize: itemSize,
		dim:      dim,
		pool:     pool,
	}
}

func (o *testNativeRec) Add(str string) {
	o.data = strings.Join([]string{o.data, str}, ";")
}

func (o *testNativeRec) Dim() []uint {
	return o.dim
}

func (o *testNativeRec) ItemSize() uint {
	return o.itemSize
}

func (o *testNativeRec) ClearData() PoolData {
	o.Add("clear")
	return o
}

func (o *testNativeRec) FreeData() {
	o.Add("free")
	o.free = true
}

func TestPool(t *testing.T) {
	var (
		dim           = []uint{16}
		itemSize uint = 4
	)

	NewPool(func(pool IPool) interface{} {
		return newTestNativeRecExt(pool, "", itemSize, dim...)
	})
}

func TestPool_GetPut(t *testing.T) {
	var (
		dim           = []uint{16}
		itemSize uint = 4
	)

	pool := NewPool(func(pool IPool) interface{} {
		return newTestNativeRecExt(pool, "create", itemSize, dim...)
	})

	func() {
		rec := (pool.Get()).(*testNativeRec)
		defer pool.Put(rec)

		rec2 := (pool.Get()).(*testNativeRec)
		defer pool.Put(rec2)

		rec.Add("processed")
		rec2.Add("processed2")
	}()

	exist1 := (pool.Get()).(*testNativeRec).data
	exist2 := (pool.Get()).(*testNativeRec).data

	expected1 := strings.Join([]string{"create", "processed", "clear"}, ";")
	expected2 := strings.Join([]string{"create", "processed2", "clear"}, ";")

	if exist1 != expected1 && exist1 != expected2 {
		t.Error("an object life cicle is different expected. 1: ", exist1, " != ", expected1, " or ", expected2)
	}

	if exist2 != expected1 && exist2 != expected2 {
		t.Error("an object life cicle is different expected. 2: ", exist2, " != ", expected1, " or ", expected2)
	}

	exist3 := (pool.Get()).(*testNativeRec).data
	expected3 := "create"

	if exist3 != expected3 {
		t.Error("an object life cicle is different expected. 3: ", exist3, " != ", expected3)
	}
}

func TestPool_FreeData(t *testing.T) {
	var (
		dim              = []uint{4}
		itemSize    uint = 16
		recordCount      = 10
	)

	pool := NewPool(func(pool IPool) interface{} {
		return newTestNativeRecExt(pool, "create", itemSize, dim...)
	})

	var records []*testNativeRec

	for range Range(0, recordCount) {
		records = append(records, func() *testNativeRec {
			rec := pool.Get().(*testNativeRec)
			defer pool.Put(rec)

			return rec
		}())
	}

	pool.FreeData()

	// check arraypool after free
	func() {
		pool.Put(newTestNativeRecExt(pool, "create", itemSize, dim...))

		if pool.Get() != nil {
			t.Error("get data after free data")
		}

		pool.FreeData()
	}()

	existCount := 0
	for _, rec := range records {
		if rec.free {
			existCount++
		}
	}

	if existCount != recordCount {
		t.Error("failed to free all records. Freed is ", existCount, ", but expected is ", recordCount)
	}
}

func TestPoolManager(t *testing.T) {
	NewPoolManager(1, func(IPool, ...uint) interface{} { return 5 })
}

func TestPoolManager_GetPut(t *testing.T) {
	var (
		dim            = []uint{16}
		itemSize1 uint = 16
		itemSize2 uint = 32
	)
	manager1 := NewPoolManager(itemSize1, func(pool IPool, dim ... uint) interface{} {
		return newTestNativeRecExt(pool, testKeyData(itemSize1, dim[0]), itemSize1, dim...)
	})
	manager2 := NewPoolManager(itemSize2, func(pool IPool, dim ... uint) interface{} {
		return newTestNativeRecExt(pool, testKeyData(itemSize2, dim[0]), itemSize2, dim...)
	})

	rec := manager1.Get(dim...).(*testNativeRec)
	if err := manager1.Put(rec); err != nil {
		t.Error("failed to put into arraypool right record with", err)
	}

	rec2 := manager2.Get(dim...).(*testNativeRec)
	if err := manager1.Put(rec2); err == nil {
		t.Error("failed to check record item size and got no error")
	}
}

func testKeyData(itemSize, len uint) string {
	return fmt.Sprintf("Key%d_%d:create", itemSize, len)
}

func testPoolInit(manager *PoolManager, prefix string, itemCount uint, itemSize uint, dims ... []uint) []*testNativeRec {
	var records []*testNativeRec

	for k := 0; k < int(itemCount); k++ {
		for i, dim := range dims {
			func() {
				for j := 1; j <= int(itemCount); j++ {
					rec := (manager.Get(dim...)).(*testNativeRec)
					defer manager.Put(rec)

					rec.Add(fmt.Sprint(prefix, i+1, j))

					if k == 0 {
						records = append(records, rec)
					}
				}
			}()
		}
	}

	return records
}

func generateExpected(prefix string, itemCount uint, itemSize uint, dims ... []uint) (data map[string]bool) {
	data = make(map[string]bool)

	for _, dim := range dims {
		id := testKeyData(itemSize, dim[0])

		for i := 1; i <= int(itemCount); i++ {
			key := id

			for j := 1; j <= int(itemCount); j++ {
				key += ";" + fmt.Sprint(prefix, i, j) + ";clear"
			}

			data[key] = false
		}
	}

	return
}

func TestPoolManager_GetWithVariants(t *testing.T) {
	var (
		itemSize  uint = 16
		itemCount uint = 4
		dims           = [][]uint{{1}, {2}, {4}, {8}}
		prefix         = "get"
	)
	manager := NewPoolManager(itemSize, func(pool IPool, dim ...uint) interface{} {
		return newTestNativeRecExt(pool, testKeyData(itemSize, dim[0]), itemSize, dim...)
	})

	// stage 1: generate data
	testPoolInit(manager, prefix, itemCount, itemSize, dims...)
	// stage 2: prepare expected record lifecicle data
	expected := generateExpected(prefix, itemCount, itemSize, dims...)

	expectedLen := len(dims) * int(itemCount)
	if existLen := len(expected); existLen != expectedLen {
		t.Error("failed to generate all possible objects lifecycle. Got ", existLen, " options, but expected is ", expectedLen)
	}
	// stage 3: compare data with pools
	for _, dim := range dims {
		func() {
			rec := (manager.Get(dim...)).(*testNativeRec)
			defer manager.Put(rec)

			if _, ok := expected[rec.data]; ok {
				expected[rec.data] = true
			}
		}()
	}
	// stage 4: count catching data
	existCount := MapIterator(expected).
		Filter(func(_ string, value bool) bool { return value }).
		Len()
	// stage 5: test it
	if expectedCount := len(dims); existCount != expectedCount {
		t.Error("failed to check updating record in the arraypool manager. Changed ", existCount, " record, but expected ", expectedCount)
	}
}

func TestPoolManager_FreeData(t *testing.T) {
	var (
		itemSize  uint = 16
		itemCount uint = 4
		dims           = [][]uint{{1}, {2}, {4}, {8}}
		prefix         = "get"
	)
	manager := NewPoolManager(itemSize, func(pool IPool, dim ... uint) interface{} {
		return newTestNativeRecExt(pool, testKeyData(itemSize, dim[0]), itemSize, dim...)
	})

	// generate data
	records := testPoolInit(manager, prefix, itemCount, itemSize, dims...)

	manager.FreeData()

	// check arraypool after free
	func() {
		dim := dims[1]

		manager.Put(newTestNativeRec("create", itemSize, dim...))

		if manager.Get(dim...) != nil {
			t.Error("get data after free data")
		}

		manager.FreeData()
	}()

	existCount := 0
	for _, rec := range records {
		if rec.free {
			existCount++
		}
	}

	if existCount != len(records) {
		t.Error("failed to free all records. Freed is ", existCount, ", but expected is ", len(records))
	}

}
