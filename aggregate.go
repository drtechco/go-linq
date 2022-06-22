package linq

// Aggregate applies an accumulator function over a sequence.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source
// except the first one. Each time f() is called, Aggregate passes both the
// element from the sequence and an aggregated value (as the first argument to
// f()). The first element of source is used as the initial aggregate value. The
// result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f().
func (q Query[T]) Aggregate(f func(interface{}, interface{}) interface{}) interface{} {
	next := q.Iterate()

	result, any := next()
	if !any {
		return nil
	}

	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}

	return result
}

// AggregateT is the typed version of Aggregate.
//
//   - f is of type: func(TSource, TSource) TSource
//
// NOTE: Aggregate has better performance than AggregateT.
func (q Query[T]) AggregateT(f func(T, T) T) T {
	next := q.Iterate()
	result, any := next()
	if !any {
		return nil
	}
	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}
	return result
}

// AggregateWithSeed applies an accumulator function over a sequence. The
// specified seed value is used as the initial accumulator value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source
// except the first one. Each time f() is called, Aggregate passes both the
// element from the sequence and an aggregated value (as the first argument to
// f()). The value of the seed parameter is used as the initial aggregate value.
// The result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f().
func (q Query[T]) AggregateWithSeed(seed interface{},
	f func(interface{}, interface{}) interface{}) interface{} {

	next := q.Iterate()
	result := seed

	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}

	return result
}

// AggregateWithSeedT is the typed version of AggregateWithSeed.
//
//   - f is of type "func(TAccumulate, TSource) TAccumulate"
//
// NOTE: AggregateWithSeed has better performance than
// AggregateWithSeedT.
func (q Query[T]) AggregateWithSeedT(seed T,
	f func(T, T) T) T {
	next := q.Iterate()
	result := seed
	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}
	return result
}

// AggregateWithSeedBy applies an accumulator function over a sequence. The
// specified seed value is used as the initial accumulator value, and the
// specified function is used to select the result value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source.
// Each time func is called, Aggregate passes both the element from the sequence
// and an aggregated value (as the first argument to func). The value of the
// seed parameter is used as the initial aggregate value. The result of func
// replaces the previous aggregated value.
//
// The final result of func is passed to resultSelector to obtain the final
// result of Aggregate.
func (q Query[T]) AggregateWithSeedBy(seed interface{},
	f func(interface{}, interface{}) interface{},
	resultSelector func(interface{}) interface{}) interface{} {

	next := q.Iterate()
	result := seed

	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}

	return resultSelector(result)
}

// AggregateWithSeedByT is the typed version of AggregateWithSeedBy.
//
//   - f is of type "func(TAccumulate, TSource) TAccumulate"
//   - resultSelectorFn is of type "func(TAccumulate) TResult"
//
// NOTE: AggregateWithSeedBy has better performance than
// AggregateWithSeedByT.
func (q Query[T]) AggregateWithSeedByT(seed T,
	f func(T, T) T,
	resultSelector func(T) T) T {
	next := q.Iterate()
	result := seed
	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}
	return resultSelector(result)
}
