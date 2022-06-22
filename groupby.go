package linq

// Group is a type that is used to store the result of GroupBy method.
type Group[K, V any] struct {
	Key   K
	Group []V
}

// GroupBy method groups the elements of a collection according to a specified
// key selector function and projects the elements for each group by using a
// specified function.
func (q Query[T]) GroupBy(keySelector func(interface{}) interface{},
	elementSelector func(interface{}) interface{}) Query[Group[any, any]] {
	return Query[Group[any, any]]{
		func() Iterator[Group[any, any]] {
			next := q.Iterate()
			set := make(map[interface{}][]interface{})

			for item, ok := next(); ok; item, ok = next() {
				key := keySelector(item)
				set[key] = append(set[key], elementSelector(item))
			}
			len := len(set)
			idx := 0
			groups := make([]Group[any, any], len)
			for k, v := range set {
				groups[idx] = Group[any, any]{k, v}
				idx++
			}
			index := 0
			return func() (item Group[any, any], ok bool) {
				ok = index < len
				if ok {
					item = groups[index]
					index++
				}

				return
			}
		},
	}
}

func GroupBy[T, M comparable, H any](q Query[T], keySelector func(T) M,
	elementSelector func(T) H) Query[Group[M, H]] {
	return Query[Group[M, H]]{
		func() Iterator[Group[M, H]] {
			next := q.Iterate()
			set := make(map[M][]H)

			for item, ok := next(); ok; item, ok = next() {
				key := keySelector(item)
				set[key] = append(set[key], elementSelector(item))
			}
			len := len(set)
			idx := 0
			groups := make([]Group[M, H], len)
			for k, v := range set {
				groups[idx] = Group[M, H]{k, v}
				idx++
			}
			index := 0
			return func() (item Group[M, H], ok bool) {
				ok = index < len
				if ok {
					item = groups[index]
					index++
				}
				return
			}
		},
	}
}
