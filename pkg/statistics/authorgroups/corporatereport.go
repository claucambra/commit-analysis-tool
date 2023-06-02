package authorgroups

import "github.com/claucambra/commit-analysis-tool/internal/db"

type CorporateReport struct {
	CorporateGroup *GroupData
	CommunityGroup *GroupData

	// Correlations based upon year-by-year aggregated figures for both groups
	InsertionsCorrel float64
	DeletionsCorrel  float64
	AuthorsCorrel    float64

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

	insertionsCorrel, deletionsCorrel, authorsCorrel := corpGroup.Correlation(commGroup)

	return &CorporateReport{
		CorporateGroup: corpGroup,
		CommunityGroup: commGroup,

		InsertionsCorrel: insertionsCorrel,
		DeletionsCorrel:  deletionsCorrel,
		AuthorsCorrel:    authorsCorrel,

		domainGroupsReport: domainGroupsReport,
	}
}
