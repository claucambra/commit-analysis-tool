package authorgroups

import (
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

// Data for an entire group
type GroupData struct {
	GroupName string

	Authors     common.EmailSet
	LineChanges *common.LineChanges
	Commits     common.CommitMap

	AuthorsPercent    float64
	InsertionsPercent float64
	DeletionsPercent  float64
}

func NewGroupData(report *DomainGroupsReport,
	groupName string,
	groupAuthors common.EmailSet,
	groupLineChanges *common.LineChanges,
	groupCommits common.CommitMap) *GroupData {

	groupData := new(GroupData)
	groupData.GroupName = groupName
	groupData.Authors = groupAuthors
	groupData.LineChanges = groupLineChanges
	groupData.Commits = groupCommits
	groupData.AuthorsPercent = (float64(len(groupAuthors)) / float64(len(report.TotalAuthors))) * 100
	groupData.InsertionsPercent = (float64(groupLineChanges.NumInsertions) / float64(report.TotalChanges.NumInsertions)) * 100
	groupData.DeletionsPercent = (float64(groupLineChanges.NumDeletions) / float64(report.TotalChanges.NumDeletions)) * 100

	return groupData
}
