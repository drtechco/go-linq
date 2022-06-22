package linq

// GroupJoin correlates the elements of two collections based on key equality,
// and groups the results.
//
// This method produces hierarchical results, which means that elements from
// outer query are paired with collections of matching elements from inner.
// GroupJoin enables you to base your results on a whole set of matches for each
// element of outer query.
//
// The resultSelector function is called only one time for each outer element
// together with a collection of all the inner elements that match the outer
// element. This differs from the Join method, in which the result selector
// function is invoked on pairs that contain one element from outer and one
// element from inner.
//
// GroupJoin preserves the order of the elements of outer, and for each element
// of outer, the order of the matching elements from inner.
func (q Query[T]) GroupJoin(inner Query[any],
	outerKeySelector func(interface{}) interface{},
	innerKeySelector func(interface{}) interface{},
	resultSelector func(outer interface{}, inners []interface{}) interface{}) Query[any] {

	return Query[any]{
		Iterate: func() Iterator[any] {
			outernext := q.Iterate()
			innernext := inner.Iterate()

			innerLookup := make(map[interface{}][]interface{})
			for innerItem, ok := innernext(); ok; innerItem, ok = innernext() {
				innerKey := innerKeySelector(innerItem)
				innerLookup[innerKey] = append(innerLookup[innerKey], innerItem)
			}

			return func() (item interface{}, ok bool) {
				if item, ok = outernext(); !ok {
					return
				}

				if group, has := innerLookup[outerKeySelector(item)]; !has {
					item = resultSelector(item, []interface{}{})
				} else {
					item = resultSelector(item, group)
				}

				return
			}
		},
	}
}

// GroupJoinT is the typed version of GroupJoin.
//
//   - inner: The query to join to the outer query.
//   - outerKeySelectorFn is of type "func(TOuter) TKey"
//   - innerKeySelectorFn is of type "func(TInner) TKey"
//   - resultSelectorFn: is of type "func(TOuter, inners []TInner) TResult"
//
// NOTE: GroupJoin has better performance than GroupJoinT.
func GroupJoin[T, M any, K comparable, O any](q Query[T], inner Query[M],
	outerKeySelector func(T) K,
	innerKeySelector func(M) K,
	resultSelector func(outer T, inners []M) O) Query[O] {

	return Query[O]{
		Iterate: func() Iterator[O] {
			outernext := q.Iterate()
			innernext := inner.Iterate()

			innerLookup := make(map[K][]M)
			for innerItem, ok := innernext(); ok; innerItem, ok = innernext() {
				innerKey := innerKeySelector(innerItem)
				innerLookup[innerKey] = append(innerLookup[innerKey], innerItem)
			}
			return func() (item O, ok bool) {
				var oItem T
				if oItem, ok = outernext(); !ok {
					return
				}
				if group, has := innerLookup[outerKeySelector(oItem)]; !has {
					item = resultSelector(oItem, []M{})
				} else {
					item = resultSelector(oItem, group)
				}
				return
			}
		},
	}
}
