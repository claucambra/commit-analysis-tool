package authorgroups

import (
	"github.com/claucambra/commit-analysis-tool/internal/db"
)

type CorporateReport struct {
	CorporateGroup *GroupData
	CommunityGroup *GroupData

	domainGroupsReport *DomainGroupsReport
}

func NewCorporateReport(groupsOfDomains map[string][]string, sqlb *db.SQLiteBackend, corporateGroupName string) *CorporateReport {
	if corporateGroupName == "" {
		corporateGroupName = "Corporate"
	}

	domainGroupsReport := NewDomainGroupsReport(groupsOfDomains)
	domainGroupsReport.Generate(sqlb)

	corpGroup := domainGroupsReport.GroupData(corporateGroupName)
	commGroup := domainGroupsReport.UnknownGroupData()

	return &CorporateReport{
		CorporateGroup: corpGroup,
		CommunityGroup: commGroup,

		domainGroupsReport: domainGroupsReport,
	}
}
