package pool

import "testing"

func TestNewTPoolManager(t *testing.T) {
	manager := NewTPoolManager()
	defer manager.FreeData()
}

func TestTPoolManager_Get(t *testing.T) {
	manager := NewTPoolManager()
	defer manager.FreeData()

	array := manager.Get(35, 10)
	defer array.Free()

	var expectedSize uint = 35 * 10
	if existSize := array.Size(); existSize != expectedSize {
		t.Error("failed to allocate array with expected size. Got ", existSize, ", but expected is ", expectedSize)
	}
}
