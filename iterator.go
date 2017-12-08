package native

type State = int

const (
	NEXT State = iota + 1
	SKIP
	FINISH
)

type Container struct {
	data map[string]bool
	Next func(key string, value bool) (string, bool, State)
}

func (o *Container) Len() (count int) {
	for key, value := range o.data {
		_, _, finish := o.Next(key, value)
		switch finish {
		case SKIP:
			continue
		case FINISH:
			return
		case NEXT:
		}

		count++
	}

	return
}

func (o *Container) Collect() (result map[string]bool) {
	result = make(map[string]bool)

	for key, value := range o.data {
		key, value, finish := o.Next(key, value)
		switch finish {
		case SKIP:
			continue
		case FINISH:
			return
		case NEXT:
		}

		result[key] = value
	}

	return
}

func (o *Container) Min() (string, bool) {
	return "", false
}

func (o *Container) Max() (string, bool) {
	return "", false
}

func (o *Container) Fold(count int, f func(int, string, bool) int) int {
	return count
}

func (o *Container) All(func(int, string) bool) bool {
	return false
}

func (o *Container) Any(func(int, string) bool) bool {
	return false
}

/*
 */
func (o *Container) Filter(f func(string, bool) bool) *Container {
	next := o.Next

	o.Next = func(k string, v bool) (key string, value bool, finish State) {
		key, value, finish = next(k, v)

		if finish == NEXT && !f(key, value) {
			finish = SKIP
		}

		return
	}

	return o
}

func (o *Container) ForEach(func(string, bool) (string, bool, State)) *Container {
	return o
}

func (o *Container) Take(count int) *Container {
	next := o.Next
	counter := 0

	o.Next = func(k string, v bool) (key string, value bool, finish State) {
		key, value, finish = next(k, v)

		if finish == NEXT {
			if counter >= count {
				finish = FINISH
			}

			counter++
		}

		return
	}

	return o
}

/*
 * ????
 */
func (o *Container) Zip(func(string, bool) bool) *Container {
	return o
}

func Iterator(data map[string]bool) *Container {
	return &Container{
		data: data,
		Next: func(key string, value bool) (string, bool, State) {
			return key, value, NEXT
		},
	}
}
