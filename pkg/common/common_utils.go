package common

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func SortedMapKeys[K constraints.Ordered, V any, M ~map[K]V](inMap M) []K {
	sortedKeys := make([]K, len(inMap))

	i := 0
	for mapKey := range inMap {
		sortedKeys[i] = mapKey
		i++
	}

	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i] < sortedKeys[j]
	})

	return sortedKeys
}
