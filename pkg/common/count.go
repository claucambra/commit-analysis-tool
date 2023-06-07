package common

import (
	"time"
)

type YearMonthCount map[int]MonthCount
type MonthCount map[int]int

func (ymc *YearMonthCount) Flatten() []int {
	values := []int{}
	years := SortedMapKeys(*ymc)

	for _, year := range years {
		for month := int(time.January); month <= int(time.December); month++ {
			monthCount, ok := (*ymc)[year][month]
			if ok {
				values = append(values, monthCount)
			}
		}
	}

	return values
}
