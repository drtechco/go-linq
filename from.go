package linq

// Iterator is an alias for function to iterate over data.
type Iterator[T any] func() (item T, ok bool)

// Query is the type returned from query functions. It can be iterated manually
// as shown in the example.
type Query[T any] struct {
	Iterate func() Iterator[T]
}

// KeyValue is a type that is used to iterate over a map (if query is created
// from a map). This type is also used by ToMap() method to output result of a
// query into a map.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// Iterable is an interface that has to be implemented by a custom collection in
// order to work with linq.
type Iterable[T any] interface {
	Iterate() Iterator[T]
}

func FromSlice[T any](source []T) Query[T] {
	arrLen := len(source)
	return Query[T]{
		Iterate: func() Iterator[T] {
			index := 0
			return func() (item T, ok bool) {
				ok = index < arrLen
				if ok {
					item = source[index]
					index++
				}
				return
			}
		},
	}
}

func FromMap[K comparable, V any](source map[K]V) Query[KeyValue[K, V]] {
	mapLen := len(source)
	return Query[KeyValue[K, V]]{
		Iterate: func() Iterator[KeyValue[K, V]] {
			index := 0
			keys := make([]K, 0)
			for k, _ := range source {
				keys = append(keys, k)
			}
			return func() (item KeyValue[K, V], ok bool) {
				ok = index < mapLen
				if ok {
					key := keys[index]
					item = KeyValue[K, V]{
						Key:   key,
						Value: source[key],
					}
					index++
				}
				return
			}
		},
	}
}

// FromChannel initializes a linq query with passed channel, linq iterates over
// channel until it is closed.
func FromChannel[T any](source <-chan T) Query[T] {
	return Query[T]{
		Iterate: func() Iterator[T] {
			return func() (item T, ok bool) {
				item, ok = <-source
				return
			}
		},
	}
}

// FromString initializes a linq query with passed string, linq iterates over
// runes of string.
func FromString(source string) Query[rune] {
	runes := []rune(source)
	len := len(runes)
	return Query[rune]{
		Iterate: func() Iterator[rune] {
			index := 0
			return func() (item rune, ok bool) {
				ok = index < len
				if ok {
					item = runes[index]
					index++
				}

				return
			}
		},
	}
}

// FromIterable initializes a linq query with custom collection passed. This
// collection has to implement Iterable interface, linq iterates over items,
// that has to implement Comparable interface or be basic types.
func FromIterable[T any](source Iterable[T]) Query[T] {
	return Query[T]{
		Iterate: source.Iterate,
	}
}
