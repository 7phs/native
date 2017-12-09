package native

import (
	"fmt"
	"strings"
	"testing"
)

type testNativeRec struct {
	len      uint
	itemSize uint
	data     string

	pool *Pool
}

func newTestNativeRec(len, itemSize uint, data string) *testNativeRec {
	return newTestNativeRecExt(len, itemSize, data, nil)
}

func newTestNativeRecExt(len, itemSize uint, data string, pool *Pool) *testNativeRec {
	return &testNativeRec{
		len:      len,
		itemSize: itemSize,
		data:     data,
		pool:     pool,
	}
}

func (o *testNativeRec) Add(str string) {
	o.data = strings.Join([]string{o.data, str}, ";")
}

func (o *testNativeRec) Len() uint {
	return o.len
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
}

func TestPool(t *testing.T) {
	var itemSize uint = 16

	NewPool(itemSize, func(len uint, pool *Pool) interface{} {
		return newTestNativeRecExt(itemSize, len, "", pool)
	})
}

func TestGetPut(t *testing.T) {
	var itemSize uint = 16

	pool := NewPool(itemSize, func(len uint, pool *Pool) interface{} {
		return newTestNativeRecExt(itemSize, len, "create", pool)
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

	expected1 := strings.Join([]string{"create", "clear", "processed", "clear"}, ";")
	expected2 := strings.Join([]string{"create", "clear", "processed2", "clear"}, ";")

	if exist1 != expected1 && exist1 != expected2 {
		t.Error("an object life cicle is different expected. 1: ", exist1, " != ", expected1, " or ", expected2)
	}

	if exist2 != expected1 && exist2 != expected2 {
		t.Error("an object life cicle is different expected. 2: ", exist2, " != ", expected1, " or ", expected2)
	}

	exist3 := (pool.Get()).(*testNativeRec).data
	expected3 := "create;clear"

	if exist3 != expected3 {
		t.Error("an object life cicle is different expected. 3: ", exist3, " != ", expected3)
	}
}

func TestPoolManager(t *testing.T) {
	NewPoolManager(1, func(uint, *Pool) interface{} { return 5 })
}

func TestPoolManagerGetPut(t *testing.T) {
	var (
		itemSize1 uint = 16
		itemSize2 uint = 32
	)
	manager1 := NewPoolManager(itemSize1, func(len uint, pool *Pool) interface{} {
		return newTestNativeRecExt(len, itemSize1, testKeyData(itemSize1, len), pool)
	})
	manager2 := NewPoolManager(itemSize2, func(len uint, pool *Pool) interface{} {
		return newTestNativeRecExt(len, itemSize2, testKeyData(itemSize2, len), pool)
	})

	rec := manager1.Get(32).(*testNativeRec)
	if err := manager1.Put(rec); err != nil {
		t.Error("failed to put into pool right record with", err)
	}

	rec2 := manager2.Get(32).(*testNativeRec)
	if err := manager1.Put(rec2); err == nil {
		t.Error("failed to check record item size and got no error")
	}
}

func testKeyData(itemSize, len uint) string {
	return fmt.Sprintf("Key%d_%d:create", itemSize, len)
}

func generateExpected(itemSize uint, lens []uint, itemCount uint, prefix string) (data map[string]bool) {
	data = make(map[string]bool)

	for _, len := range lens {
		id := testKeyData(itemSize, len) + ";clear"

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

func TestPoolManagerGetWithVariants(t *testing.T) {
	var itemSize uint = 16
	manager := NewPoolManager(itemSize, func(len uint, pool *Pool) interface{} {
		return newTestNativeRecExt(len, itemSize, testKeyData(itemSize, len), pool)
	})
	var itemCount uint = 4
	lens := []uint{1, 2, 4, 8}
	prefix := "get"

	// stage 1: generate data
	for k := 0; k < int(itemCount); k++ {
		for i, len := range lens {
			func() {
				for j := 1; j <= int(itemCount); j++ {
					rec := (manager.Get(len)).(*testNativeRec)
					defer manager.Put(rec)

					rec.Add(fmt.Sprint(prefix, i+1, j))
				}
			}()
		}
	}
	// stage 2: prepare expected record lifecicle data
	expected := generateExpected(itemSize, lens, itemCount, prefix)

	expected_len := len(lens) * int(itemCount)
	if exist_len := len(expected); exist_len != expected_len {
		t.Error("failed to generate all possible objects lifecycle. Got ", exist_len, " options, but expected is ", expected_len)
	}
	// stage 3: compare data with pools
	for _, len := range lens {
		func() {
			rec := (manager.Get(len)).(*testNativeRec)
			defer manager.Put(rec)

			if _, ok := expected[rec.data]; ok {
				expected[rec.data] = true
			}
		}()
	}
	// stage 4: count catching data
	existCount := Iterator(expected).
		Filter(func(key string, item bool) bool {
			return item
		}).
		Len()
	// stage 5: test it
	if expectedCount := len(lens); existCount != expectedCount {
		t.Error("failed to check updating record in the pool manager. Changed ", existCount, " record, but expected ", expectedCount)
	}
}
