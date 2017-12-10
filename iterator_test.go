package native

import (
	"testing"
)

func TestMapContainer_FilterTake(t *testing.T) {
	data := map[string]bool{
		"hello":     false,
		"привет":    true,
		"morning":   false,
		"утро":      true,
		"winter":    false,
		"зима":      true,
		"movie":     false,
		"фильм":     true,
		"TV":        false,
		"ТВ":        true,
		"snowdrift": false,
		"сугроб":    true,
		"courtyard": false,
		"двор":      true,
	}

	expectedLen := 7
	existLen := MapIterator(data).
		Filter(func(word string, isRussian bool) bool { return isRussian }).
		Len()

	if existLen != expectedLen {
		t.Error("failed to filter russian word. Got ", existLen, " words, but expected is ", expectedLen, " words")
	}

	expectedLen = 3
	existLen = MapIterator(data).
		Filter(func(word string, isRussian bool) bool { return isRussian }).
		Take(expectedLen).
		Len()

	if existLen != expectedLen {
		t.Error("failed to filter russian word and take just ", expectedLen, " words. Got ", existLen, " words, but expected is ", expectedLen, " words")
	}
}

func TestUintContainer_FilterTakeMul(t *testing.T) {
	data := []uint{0, 3, 5, 7, 18, 30, 36, 46, 48, 52, 60}

	expectedLen := 7
	existLen := UintIterator(data).
		Filter(func(index int, number uint) bool { return number%3 == 0 }).
		Len()

	if existLen != expectedLen {
		t.Error("failed to filter number divided by 3. Got ", existLen, " words, but expected is ", expectedLen, " words")
	}

	expectedLen = 3
	existLen = UintIterator(data).
		Filter(func(index int, number uint) bool { return number%3 == 0 }).
		Take(expectedLen).
		Len()

	if existLen != expectedLen {
		t.Error("failed to filter number divided by 3 and take just ", expectedLen, " words. Got ", existLen, " words, but expected is ", expectedLen, " words")
	}

	var expectedMul uint = 167961600
	existMul := UintIterator(data).
		Filter(func(index int, number uint) bool { return number > 0 && number%3 == 0 }).
		Mul()

	if existMul != expectedMul {
		t.Error("failed to filter number divided by 3 and multiply it. Got ", existMul, ", but expected is ", expectedMul)
	}
}

func BenchmarkUintContainer_Len(b *testing.B) {
	dataLen := 1000
	data := make([]uint, 0, dataLen)
	for i := 0; i < cap(data); i++ {
		data = append(data, uint(i)*3+1)
	}

	b.Run("native iterator", func(b *testing.B) {
		count := 0

		for i:=0;i<100;i++ {
			for _, number := range data {
				if number > 0 && number%2 == 0 {
					count++
				}
			}
		}

		if count<10 {
			b.Error("failed to native calc length", count)
		}
	})

	iterator := UintIterator(data)
	b.Run("native iterator", func(b *testing.B) {
		count := 0

		for i:=0;i<100;i++ {
			count += iterator.
				Filter(func(index int, number uint) bool { return number > 0 && number%2 == 0 }).
				Len()
		}

		if count<10 {
			b.Error("failed to iterative calc length", count)
		}
	})
}
