package array

import (
	"reflect"
	"testing"
)

func TestWithTArray(t *testing.T) {
	expected := []int{0, 2, 4, 6, 8}

	d := WithTArray(expected)

	if exist := d.Marshal(); !reflect.DeepEqual(exist, expected) {
		t.Error("failed to make copy as native array. Got", exist, ", but expected is", expected)
	}
}
