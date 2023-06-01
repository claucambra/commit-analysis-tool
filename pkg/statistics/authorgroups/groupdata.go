package authorgroups

import (
	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"gonum.org/v1/gonum/stat"
)

type GroupData struct {
	GroupName string

	NumAuthors        int
	NumLineChanges    *common.LineChanges
	YearlyLineChanges common.YearlyLineChangeMap

	AuthorsPercent    float32
	InsertionsPercent float32
	DeletionsPercent  float32
}

func NewGroupData(report *DomainGroupsReport, groupName string, groupAuthors int, groupLineChanges *common.LineChanges, groupYearlyLineChanges common.YearlyLineChangeMap) *GroupData {
	groupData := new(GroupData)
	groupData.GroupName = groupName
	groupData.NumAuthors = groupAuthors
	groupData.NumLineChanges = groupLineChanges
	groupData.YearlyLineChanges = groupYearlyLineChanges
	groupData.AuthorsPercent = (float32(groupAuthors) / float32(report.TotalAuthors)) * 100
	groupData.InsertionsPercent = (float32(groupLineChanges.NumInsertions) / float32(report.TotalChanges.NumInsertions)) * 100
	groupData.DeletionsPercent = (float32(groupLineChanges.NumDeletions) / float32(report.TotalChanges.NumDeletions)) * 100

	return groupData
}

// Returns correlations between insertions and deletions of two groups
func (group *GroupData) Correlation(groupToCorrelate *GroupData) (float64, float64) {
	yearsInBothGroups := []int{}
	for year := range group.YearlyLineChanges {
		_, ok := groupToCorrelate.YearlyLineChanges[year]
		if ok {
			yearsInBothGroups = append(yearsInBothGroups, year)
		}
	}

	thisYearlyInsertionsArr, thisYearlyDeletionsArr := group.YearlyLineChanges.SeparatedChangeArrays(yearsInBothGroups)
	thatYearlyInsertionsArr, thatYearlyDeletionsArr := groupToCorrelate.YearlyLineChanges.SeparatedChangeArrays(yearsInBothGroups)

	floatThisInserts := make([]float64, len(yearsInBothGroups))
	floatThisDeletes := make([]float64, len(yearsInBothGroups))
	floatThatInserts := make([]float64, len(yearsInBothGroups))
	floatThatDeletes := make([]float64, len(yearsInBothGroups))

	for i := 0; i < len(group.YearlyLineChanges); i++ {
		floatThisInserts[i] = float64(thisYearlyInsertionsArr[i])
		floatThisDeletes[i] = float64(thisYearlyDeletionsArr[i])
	}
	for i := 0; i < len(groupToCorrelate.YearlyLineChanges); i++ {
		floatThatInserts[i] = float64(thatYearlyInsertionsArr[i])
		floatThatDeletes[i] = float64(thatYearlyDeletionsArr[i])
	}

	insertionsCorrel := stat.Correlation(floatThisInserts, floatThatInserts, nil)
	deletionsCorrel := stat.Correlation(floatThisDeletes, floatThatDeletes, nil)

	return insertionsCorrel, deletionsCorrel
}
