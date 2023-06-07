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

func MinInt(numA int, numB int) int {
	minNum := math.Min(float64(numA), float64(numB))
	return int(minNum)
}

func SliceContains[V comparable](slice []V, valueToFind V) (bool, int) {
	for i, value := range slice {
		if value == valueToFind {
			return true, i
		}
	}

	return false, -1
}

func SliceIntToFloat[I constraints.Integer, F constraints.Float](slice []I) []F {
	outSlice := make([]F, len(slice))

	for i, intVal := range slice {
		outSlice[i] = F(intVal)
	}

	return outSlice
}

// Returns whichever map has the higher key at the beginning when sorted
func HigherStartKey[K constraints.Ordered, V any, M ~map[K]V](inMapA M, inMapB M) K {
	sortedAKeys := SortedMapKeys(inMapA)
	sortedBKeys := SortedMapKeys(inMapB)

	sortedAKeysLen := len(sortedAKeys)
	sortedBKeysLen := len(sortedBKeys)

	if sortedAKeysLen == 0 {
		return sortedBKeys[0]
	} else if sortedBKeysLen == 0 {
		return sortedAKeys[0]
	} else {
		firstSortedAKey := sortedAKeys[0]
		firstSortedBKey := sortedBKeys[0]

		if sortedAKeys[0] > sortedBKeys[0] {
			return firstSortedAKey
		} else {
			return firstSortedBKey
		}
	}
}
