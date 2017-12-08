package native

import (
	"fmt"
	"strings"
	"testing"
)

type testNativeRec struct {
	data string

	pool *Pool
}

func newTestNativeRec(data string) *testNativeRec {
	return newTestNativeRecExt(data, nil)
}

func newTestNativeRecExt(data string, pool *Pool) *testNativeRec {
	return &testNativeRec{
		data: data,
		pool: pool,
	}
}

func (o *testNativeRec) Add(str string) {
	o.data = strings.Join([]string{o.data, str}, ";")
}

func (o *testNativeRec) ClearData() PoolData {
	o.Add("clear")
	return o
}

func (o *testNativeRec) FreeData() {
	o.Add("free")
}

func TestPool(t *testing.T) {
	NewPool(func(pool *Pool) interface{} {
		return newTestNativeRecExt("", pool)
	})
}

func TestGetPut(t *testing.T) {
	pool := NewPool(func(pool *Pool) interface{} {
		return newTestNativeRecExt("create", pool)
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
	NewPoolManager()
}

func testKeyData(kind, itemSize uint) string {
	return fmt.Sprintf("Key%d_%d:create", kind, itemSize)
}

func newTestKey(kind, itemSize uint) *PoolId {
	return &PoolId{
		Kind:     kind,
		ItemSize: itemSize,
		New: func(pool *Pool) interface{} {
			return newTestNativeRecExt(testKeyData(kind, itemSize), pool)
		},
	}
}

func generateTestRec(manager *PoolManager, key *PoolId) (rec1, rec2 *testNativeRec) {
	rec1 = (manager.Get(key)).(*testNativeRec)
	rec2 = (manager.Get(key)).(*testNativeRec)

	return rec1, rec2
}

func generateExpected(keys []*PoolId, part1 []string, part2 []string) (data map[string]bool) {
	data = make(map[string]bool)

	for _, key := range keys {
		for _, str1 := range part1 {
			for _, str2 := range part2 {
				id := testKeyData(key.Kind, key.ItemSize) + ";clear;" + str1 + ";clear;" + str2 + ";clear"
				data[id] = false
			}
		}
	}

	return
}

func TestPoolManagerGet(t *testing.T) {
	manager := NewPoolManager()

	// prepare key
	Key1_16 := newTestKey(1, 16)
	Key1_32 := newTestKey(1, 32)
	Key2_16 := newTestKey(2, 16)
	Key2_32 := newTestKey(2, 32)

	keyList := []*PoolId{
		Key1_16, Key1_32, Key2_16, Key2_32,
	}
	// stage 1: generate data
	for _, key := range keyList {
		func(rec1, rec2 *testNativeRec) {
			defer manager.Put(key, rec1)
			defer manager.Put(key, rec2)

			rec1.Add("get11")
			rec2.Add("get12")
		}(generateTestRec(manager, key))
	}
	// stage 2: generate data
	for _, key := range keyList {
		func(rec1, rec2 *testNativeRec) {
			defer manager.Put(key, rec1)
			defer manager.Put(key, rec2)

			rec1.Add("get21")
			rec2.Add("get22")
		}(generateTestRec(manager, key))
	}
	// prepare expeted data
	expected := generateExpected(keyList,
		[]string{"get11", "get12"},
		[]string{"get21", "get22"})
	// compare data with pools
	for _, key := range keyList {
		func() {
			rec := (manager.Get(key)).(*testNativeRec)
			defer manager.Put(key, rec)

			for id := range expected {
				if rec.data == id {
					expected[id] = true
					break
				}
			}
		}()
	}
	// count catching data
	existCount := Iterator(expected).
		Filter(func(key string, item bool) bool {
			return item
		}).
		Len()
	// test it
	if expectedCount := len(keyList); existCount != expectedCount {
		t.Error("failed to check updating record in the pool manager. Changed ", existCount, " record, but expected ", expectedCount)
	}
}
