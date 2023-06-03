package authorgroups

import (
	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"gonum.org/v1/gonum/stat"
)

// Data for an entire group
type GroupData struct {
	GroupName string

	Authors           common.EmailSet
	LineChanges       *common.LineChanges
	YearlyLineChanges common.YearlyLineChangeMap
	YearlyAuthors     common.YearlyEmailMap

	AuthorsPercent    float32
	InsertionsPercent float32
	DeletionsPercent  float32
}

func NewGroupData(report *DomainGroupsReport,
	groupName string,
	groupAuthors common.EmailSet,
	groupLineChanges *common.LineChanges,
	groupYearlyLineChanges common.YearlyLineChangeMap,
	groupYearlyAuthors common.YearlyEmailMap) *GroupData {

	groupData := new(GroupData)
	groupData.GroupName = groupName
	groupData.Authors = groupAuthors
	groupData.LineChanges = groupLineChanges
	groupData.YearlyLineChanges = groupYearlyLineChanges
	groupData.YearlyAuthors = groupYearlyAuthors
	groupData.AuthorsPercent = (float32(len(groupAuthors)) / float32(len(report.TotalAuthors))) * 100
	groupData.InsertionsPercent = (float32(groupLineChanges.NumInsertions) / float32(report.TotalChanges.NumInsertions)) * 100
	groupData.DeletionsPercent = (float32(groupLineChanges.NumDeletions) / float32(report.TotalChanges.NumDeletions)) * 100

	return groupData
}

func (group *GroupData) changesCorrelation(groupToCorrelate *GroupData) (float64, float64) {
	changeYearsInBothGroups := common.KeysInCommon(group.YearlyLineChanges, groupToCorrelate.YearlyLineChanges)

	thisYearlyInsertionsArr, thisYearlyDeletionsArr := group.YearlyLineChanges.SeparatedChangeArrays(changeYearsInBothGroups)
	thatYearlyInsertionsArr, thatYearlyDeletionsArr := groupToCorrelate.YearlyLineChanges.SeparatedChangeArrays(changeYearsInBothGroups)

	floatThisInserts := make([]float64, len(changeYearsInBothGroups))
	floatThisDeletes := make([]float64, len(changeYearsInBothGroups))
	floatThatInserts := make([]float64, len(changeYearsInBothGroups))
	floatThatDeletes := make([]float64, len(changeYearsInBothGroups))

	for i := 0; i < len(changeYearsInBothGroups); i++ {
		floatThisInserts[i] = float64(thisYearlyInsertionsArr[i])
		floatThisDeletes[i] = float64(thisYearlyDeletionsArr[i])
		floatThatInserts[i] = float64(thatYearlyInsertionsArr[i])
		floatThatDeletes[i] = float64(thatYearlyDeletionsArr[i])
	}

	insertionsCorrel := stat.Correlation(floatThisInserts, floatThatInserts, nil)
	deletionsCorrel := stat.Correlation(floatThisDeletes, floatThatDeletes, nil)

	return insertionsCorrel, deletionsCorrel
}

func (group *GroupData) authorsCorrelation(groupToCorrelate *GroupData) float64 {
	authorYearsInBothGroups := common.KeysInCommon(group.YearlyAuthors, groupToCorrelate.YearlyAuthors)

	thisYearlyAuthorsArr := group.YearlyAuthors.CountArray(authorYearsInBothGroups)
	thatYearlyAuthorsArr := groupToCorrelate.YearlyAuthors.CountArray(authorYearsInBothGroups)

	floatThisAuthors := make([]float64, len(authorYearsInBothGroups))
	floatThatAuthors := make([]float64, len(authorYearsInBothGroups))

	for i := 0; i < len(authorYearsInBothGroups); i++ {
		floatThisAuthors[i] = float64(thisYearlyAuthorsArr[i])
		floatThatAuthors[i] = float64(thatYearlyAuthorsArr[i])
	}

	return stat.Correlation(floatThisAuthors, floatThatAuthors, nil)
}

// Returns correlations between insertions and deletions and authors of two groups over years
func (group *GroupData) Correlation(groupToCorrelate *GroupData) (float64, float64, float64) {
	insertionsCorrel, deletionsCorrel := group.changesCorrelation(groupToCorrelate)
	authorsCorrel := group.authorsCorrelation(groupToCorrelate)
	return insertionsCorrel, deletionsCorrel, authorsCorrel
}
