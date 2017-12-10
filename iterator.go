package native

type State = int

const (
	NEXT State = iota + 1
	SKIP
	FINISH
)

type MapContainer struct {
	data map[string]bool

	Next func(index string, value bool) (string, bool, State)
}

func MapIterator(data map[string]bool) *MapContainer {
	return &MapContainer{
		data: data,
		Next: func(index string, value bool) (string, bool, State) {
			return index, value, NEXT
		},
	}
}

func (o *MapContainer) Filter(f func(string, bool) bool) *MapContainer {
	next := o.Next

	o.Next = func(i string, v bool) (index string, value bool, state State) {
		if index, value, state = next(i, v); state == NEXT && !f(index, value) {
			state = SKIP
		}

		return
	}

	return o
}

func (o *MapContainer) Take(count int) *MapContainer {
	next := o.Next
	counter := 0

	o.Next = func(i string, v bool) (index string, value bool, state State) {
		if index, value, state = next(i, v); state == NEXT {
			if counter >= count {
				state = FINISH
			}

			counter++
		}

		return
	}

	return o
}

func (o *MapContainer) Range(fn func(index string, value bool)) {
	for index, value := range o.data {
		switch index, value, status := o.Next(index, value); status {
		case SKIP:
			continue
		case FINISH:
			return
		case NEXT:
			fn(index, value)
		}
	}
}

func (o *MapContainer) Len() (count int) {
	o.Range(func(_ string, _ bool) {
		count++
	})

	return
}

type UintContainer struct {
	data []uint

	Next func(index int, value uint) (int, uint, State)
}

func UintIterator(data []uint) *UintContainer {
	return &UintContainer{
		data: data,
		Next: func(index int, value uint) (int, uint, State) {
			return index, value, NEXT
		},
	}
}

func (o *UintContainer) Filter(f func(int, uint) bool) *UintContainer {
	next := o.Next

	o.Next = func(i int, v uint) (index int, value uint, state State) {
		if index, value, state = next(i, v); state == NEXT && !f(index, value) {
			state = SKIP
		}

		return
	}

	return o
}

func (o *UintContainer) Take(count int) *UintContainer {
	next := o.Next
	counter := 0

	o.Next = func(i int, v uint) (index int, value uint, state State) {
		if index, value, state = next(i, v); state == NEXT {
			if counter >= count {
				state = FINISH
			}

			counter++
		}

		return
	}

	return o
}

func (o *UintContainer) Range(fn func(index int, value uint)) {
	for index, value := range o.data {
		switch index, value, status := o.Next(index, value); status {
		case SKIP:
			continue
		case FINISH:
			return
		case NEXT:
			fn(index, value)
		}
	}
}

func (o *UintContainer) Mul() (mul uint) {
	mul = 1

	o.Range(func(index int, value uint) {
		mul *= value
	})

	return
}

func (o *UintContainer) Len() (count int) {
	o.Range(func(_ int, _ uint) {
		count++
	})

	return
}

func Range(begin, end int) <-chan int {
	yield := make(chan int)

	go func() {
		for ; begin < end; begin++ {
			yield <- begin
		}

		close(yield)
	}()

	return yield
}
