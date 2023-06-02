package common

import (
	"sort"
)

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

func (ypm *YearlyPeopleMap) AddYearlyPeopleMap(ypmToAdd YearlyPeopleMap) {
	for year, peopleToAdd := range ypmToAdd {
		if existingPeople, ok := (*ypm)[year]; ok {
			(*ypm)[year] = append(existingPeople, peopleToAdd...)
		} else {
			(*ypm)[year] = peopleToAdd
		}
	}
}

func (ypm *YearlyPeopleMap) SubtractYearlyPeopleMap(ypmToSubtract YearlyPeopleMap) {
	for year, peopleToSubtract := range ypmToSubtract {
		if existingPeople, ok := (*ypm)[year]; ok {
			subtractedPeople := existingPeople

			for _, personToSubtract := range peopleToSubtract {
				personIdx := sort.Search(len(subtractedPeople), func(i int) bool {
					return subtractedPeople[i].Email == personToSubtract.Email
				})

				if personIdx < len(subtractedPeople) && subtractedPeople[personIdx].Email == personToSubtract.Email {
					subtractedPeople[personIdx] = subtractedPeople[len(subtractedPeople)-1]
					subtractedPeople = subtractedPeople[:len(subtractedPeople)-1]
				}
			}

			(*ypm)[year] = subtractedPeople
		}
	}
}
