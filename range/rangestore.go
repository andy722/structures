package _range

import "sort"

type RangePoint = uint64

// RangeStore maps inclusive range [fromIncl, toIncl] to value
type RangeStore struct {
	ranges []RangePoint
	values []RangeStoreItem
}

type RangeStoreItem struct {
	toIncl RangePoint
	value  interface{}
}

func NewRangeStore() RangeStore {
	return RangeStore{
		ranges: make([]RangePoint, 0, 10000),
		values: make([]RangeStoreItem, 0, 10000),
	}
}

// Add associates range of [fromIncl, toIncl] to value
func (store *RangeStore) Add(fromIncl RangePoint, toIncl RangePoint, value interface{}) {
	store.ranges = append(store.ranges, fromIncl)

	item := RangeStoreItem{
		toIncl: toIncl,
		value:  value,
	}
	store.values = append(store.values, item)
}

// Lookup returns a value associated with a range including specified point, or nil if none found
func (store *RangeStore) Lookup(point RangePoint) interface{} {
	count := len(store.ranges)

	i := sort.Search(count, func(i int) bool {
		return store.ranges[i] >= point
	})

	if i == count {
		return nil
	}

	val := store.checkMatch(i, point)
	if val != nil {
		return val
	}

	if i > 0 {
		val := store.checkMatch(i-1, point)
		if val != nil {
			return val
		}
	}
	return nil
}

func (store *RangeStore) checkMatch(i int, value RangePoint) interface{} {
	if store.ranges[i] > value {
		return nil
	}

	if found := &store.values[i]; found.toIncl >= value {
		return found.value
	}

	return nil
}
