package common

type Person struct {
	Name  string
	Email string
}

type YearlyPeopleMap map[int][]*Person

func (ypm *YearlyPeopleMap) CountArray(years []int) []int {
	yearsToReturn := years

	if yearsToReturn == nil {
		yearsToReturn = SortedMapKeys(*ypm)
	}

	countArray := make([]int, len(yearsToReturn))

	for i, year := range yearsToReturn {
		countArray[i] = len((*ypm)[year])
	}

	return countArray
}
