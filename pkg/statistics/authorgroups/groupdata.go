package authorgroups

import "github.com/claucambra/commit-analysis-tool/pkg/common"

type GroupData struct {
	GroupName string

	NumAuthors int
	NumChanges *common.LineChanges

	AuthorsPercent    float32
	InsertionsPercent float32
	DeletionsPercent  float32
}

func NewGroupData(report *DomainGroupsReport, groupName string, groupAuthors int, groupChanges *common.LineChanges) *GroupData {
	groupData := new(GroupData)
	groupData.GroupName = groupName
	groupData.NumAuthors = groupAuthors
	groupData.NumChanges = groupChanges
	groupData.AuthorsPercent = (float32(groupAuthors) / float32(report.TotalAuthors)) * 100
	groupData.InsertionsPercent = (float32(groupChanges.NumInsertions) / float32(report.TotalChanges.NumInsertions)) * 100
	groupData.DeletionsPercent = (float32(groupChanges.NumDeletions) / float32(report.TotalChanges.NumDeletions)) * 100

	return groupData
}
