package common

import (
	"math"
	"sort"
	"time"

	"gonum.org/v1/gonum/stat"
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

func CorrelateYearMonthCounts(ymc1 YearMonthCount, ymc2 YearMonthCount) float64 {
	commonYears := KeysInCommon(ymc1, ymc2)

	if len(commonYears) == 0 {
		return math.NaN()
	}

	sort.Slice(commonYears, func(i, j int) bool {
		return commonYears[i] < commonYears[j]
	})

	firstYear := commonYears[0] // Start from first year in common
	lastYear := commonYears[len(commonYears)-1]

	filteredYmc1 := YearMonthCount{}
	filteredYmc2 := YearMonthCount{}

	for i := firstYear; i <= lastYear; i++ {
		ymc1MonthMap, ymc1MonthMapOk := ymc1[i]
		ymc2MonthMap, ymc2MonthMapOk := ymc2[i]

		firstMonth := int(time.January)
		lastMonth := int(time.December)
		if i == firstYear {
			firstMonth = HigherStartKey(ymc1MonthMap, ymc2MonthMap)
		} else if i == lastYear {
			lastMonth = LowerEndKey(ymc1MonthMap, ymc2MonthMap)
		}

		// We fill in any gaps in the YearMonthCounts relative to each other with 0s
		for j := firstMonth; j <= lastMonth; j++ {
			// TODO: Deduplicate
			if _, ok := filteredYmc1[i]; !ok {
				filteredYmc1[i] = MonthCount{}
			}
			filteredYmc1[i][j] = 0
			if ymc1MonthMapOk {
				ymc1MonthValue, ymc1MonthValueOk := ymc1MonthMap[j]
				if ymc1MonthValueOk {
					filteredYmc1[i][j] = ymc1MonthValue
				}
			}

			if _, ok := filteredYmc2[i]; !ok {
				filteredYmc2[i] = MonthCount{}
			}
			filteredYmc2[i][j] = 0
			if ymc2MonthMapOk {
				ymc2MonthValue, ymc2MonthValueOk := ymc2MonthMap[j]
				if ymc2MonthValueOk {
					filteredYmc2[i][j] = ymc2MonthValue
				}
			}
		}
	}

	flatFloatFilteredYmc1 := SliceIntToFloat[int, float64](filteredYmc1.Flatten())
	flatFloatFilteredYmc2 := SliceIntToFloat[int, float64](filteredYmc2.Flatten())

	return stat.Correlation(flatFloatFilteredYmc1, flatFloatFilteredYmc2, nil)
}
