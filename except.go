package linq

// Except produces the set difference of two sequences. The set difference is
// the members of the first sequence that don't appear in the second sequence.
func (q Query[T]) Except(q2 Query[T]) Query[T] {
	return Query[T]{
		Iterate: func() Iterator[T] {
			next := q.Iterate()

			next2 := q2.Iterate()
			set := make(map[any]bool)
			for i, ok := next2(); ok; i, ok = next2() {
				set[i] = true
			}

			return func() (item T, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if _, has := set[item]; !has {
						return
					}
				}

				return
			}
		},
	}
}

// ExceptBy invokes a transform function on each element of a collection and
// produces the set difference of two sequences. The set difference is the
// members of the first sequence that don't appear in the second sequence.
func (q Query[T]) ExceptBy(q2 Query[T],
	selector func(interface{}) interface{}) Query[T] {
	return Query[T]{
		Iterate: func() Iterator[T] {
			next := q.Iterate()
			next2 := q2.Iterate()
			set := make(map[any]bool)
			for i, ok := next2(); ok; i, ok = next2() {
				s := selector(i)
				set[s] = true
			}
			return func() (item T, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					s := selector(item)
					if _, has := set[s]; !has {
						return
					}
				}

				return
			}
		},
	}
}

// ExceptByT is the typed version of ExceptBy.
//
//   - selectorFn is of type "func(TSource) TSource"
//
// NOTE: ExceptBy has better performance than ExceptByT.
func (q Query[T]) ExceptByT(q2 Query[T],
	selector func[M any](T) M) Query[M] {
	selectorGenericFunc, err := newGenericFunc(
		"ExceptByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.ExceptBy(q2, selectorFunc)
}
