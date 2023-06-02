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
	YearlyAuthors     common.YearlyEmailMap

	AuthorsPercent    float32
	InsertionsPercent float32
	DeletionsPercent  float32
}

func NewGroupData(report *DomainGroupsReport, groupName string, groupAuthors int, groupLineChanges *common.LineChanges, groupYearlyLineChanges common.YearlyLineChangeMap, groupYearlyAuthors common.YearlyEmailMap) *GroupData {
	groupData := new(GroupData)
	groupData.GroupName = groupName
	groupData.NumAuthors = groupAuthors
	groupData.NumLineChanges = groupLineChanges
	groupData.YearlyLineChanges = groupYearlyLineChanges
	groupData.YearlyAuthors = groupYearlyAuthors
	groupData.AuthorsPercent = (float32(groupAuthors) / float32(report.TotalAuthors)) * 100
	groupData.InsertionsPercent = (float32(groupLineChanges.NumInsertions) / float32(report.TotalChanges.NumInsertions)) * 100
	groupData.DeletionsPercent = (float32(groupLineChanges.NumDeletions) / float32(report.TotalChanges.NumDeletions)) * 100

	return groupData
}

// Returns correlations between insertions and deletions and authors of two groups over years
func (group *GroupData) Correlation(groupToCorrelate *GroupData) (float64, float64, float64) {
	changeYearsInBothGroups := []int{}
	for year := range group.YearlyLineChanges {
		_, ok := groupToCorrelate.YearlyLineChanges[year]
		if ok {
			changeYearsInBothGroups = append(changeYearsInBothGroups, year)
		}
	}

	authorYearsInBothGroups := []int{}
	for year := range group.YearlyAuthors {
		_, ok := groupToCorrelate.YearlyAuthors[year]
		if ok {
			authorYearsInBothGroups = append(authorYearsInBothGroups, year)
		}
	}

	thisYearlyInsertionsArr, thisYearlyDeletionsArr := group.YearlyLineChanges.SeparatedChangeArrays(changeYearsInBothGroups)
	thisYearlyAuthorsArr := group.YearlyAuthors.CountArray(authorYearsInBothGroups)
	thatYearlyInsertionsArr, thatYearlyDeletionsArr := groupToCorrelate.YearlyLineChanges.SeparatedChangeArrays(changeYearsInBothGroups)
	thatYearlyAuthorsArr := groupToCorrelate.YearlyAuthors.CountArray(authorYearsInBothGroups)

	floatThisInserts := make([]float64, len(changeYearsInBothGroups))
	floatThisDeletes := make([]float64, len(changeYearsInBothGroups))
	floatThisAuthors := make([]float64, len(authorYearsInBothGroups))
	floatThatInserts := make([]float64, len(changeYearsInBothGroups))
	floatThatDeletes := make([]float64, len(changeYearsInBothGroups))
	floatThatAuthors := make([]float64, len(authorYearsInBothGroups))

	for i := 0; i < len(changeYearsInBothGroups); i++ {
		floatThisInserts[i] = float64(thisYearlyInsertionsArr[i])
		floatThisDeletes[i] = float64(thisYearlyDeletionsArr[i])
		floatThatInserts[i] = float64(thatYearlyInsertionsArr[i])
		floatThatDeletes[i] = float64(thatYearlyDeletionsArr[i])
	}

	for i := 0; i < len(authorYearsInBothGroups); i++ {
		floatThisAuthors[i] = float64(thisYearlyAuthorsArr[i])
		floatThatAuthors[i] = float64(thatYearlyAuthorsArr[i])
	}

	insertionsCorrel := stat.Correlation(floatThisInserts, floatThatInserts, nil)
	deletionsCorrel := stat.Correlation(floatThisDeletes, floatThatDeletes, nil)
	authorsCorrel := stat.Correlation(floatThisAuthors, floatThatAuthors, nil)

	return insertionsCorrel, deletionsCorrel, authorsCorrel
}
