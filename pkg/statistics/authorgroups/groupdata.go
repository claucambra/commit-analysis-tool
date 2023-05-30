package authorgroups

import "github.com/claucambra/commit-analysis-tool/pkg/common"

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
	groupData.AuthorsPercent = (float32(groupAuthors) / float32(report.TotalAuthors)) * 100
	groupData.InsertionsPercent = (float32(groupLineChanges.NumInsertions) / float32(report.TotalChanges.NumInsertions)) * 100
	groupData.DeletionsPercent = (float32(groupLineChanges.NumDeletions) / float32(report.TotalChanges.NumDeletions)) * 100

	return groupData
}
