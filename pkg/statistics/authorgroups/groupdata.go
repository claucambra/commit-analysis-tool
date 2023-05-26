package authorgroups

type GroupData struct {
	GroupName string

	NumAuthors    int
	NumInsertions int
	NumDeletions  int

	AuthorsPercent    float32
	InsertionsPercent float32
	DeletionsPercent  float32
}

func NewGroupData(report *DomainGroupsReport, groupName string, groupAuthors int, groupInsertions int, groupDeletions int) *GroupData {
	groupData := new(GroupData)
	groupData.GroupName = groupName
	groupData.NumAuthors = groupAuthors
	groupData.NumInsertions = groupInsertions
	groupData.NumDeletions = groupDeletions
	groupData.AuthorsPercent = (float32(groupAuthors) / float32(report.TotalAuthors)) * 100
	groupData.InsertionsPercent = (float32(groupInsertions) / float32(report.TotalInsertions)) * 100
	groupData.DeletionsPercent = (float32(groupDeletions) / float32(report.TotalDeletions)) * 100

	return groupData
}
