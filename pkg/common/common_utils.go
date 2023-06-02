package common

import (
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

func CopyMap[K comparable, V any, M ~map[K]V](inMap M) map[K]V {
	copyMap := make(map[K]V, len(inMap))

	for key, value := range inMap {
		copyMap[key] = value
	}

	return copyMap
}

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

func KeysInCommon[K comparable, V any, M ~map[K]V](inMapA M, inMapB M) []K {
	commonKeys := []K{}

	for key := range inMapA {
		_, ok := inMapB[key]
		if ok {
			commonKeys = append(commonKeys, key)
		}
	}

	return commonKeys
}

// Insert operation that adds to an existing value in the map using a specified additive function
func AdditiveValueMapInsert[K comparable, V any, M ~map[K]V](inMap map[K]V, key K, additiveFunc func(V, V) V, valueToAdd V) {
	completeValue := valueToAdd
	if existingValue, ok := inMap[key]; ok {
		completeValue = additiveFunc(valueToAdd, existingValue)
	}

	inMap[key] = completeValue
}

// Operation that subtracts from an existing value in the map using a specified subtractive function
func SubtractiveValueMapRemove[K comparable, V any, M ~map[K]V](inMap map[K]V, key K, subtractiveFunc func(V, V) (V, bool), valueToSubtract V) {
	if existingValue, ok := inMap[key]; ok {
		subtractedValue, removeCondition := subtractiveFunc(existingValue, valueToSubtract)

		if removeCondition {
			delete(inMap, key)
		} else {
			inMap[key] = subtractedValue
		}
	}
}

func MaxInt(numA int, numB int) int {
	maxNum := math.Max(float64(numA), float64(numB))
	return int(maxNum)
}
