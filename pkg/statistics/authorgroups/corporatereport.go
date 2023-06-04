package authorgroups

import (
	"github.com/claucambra/commit-analysis-tool/internal/db"
)

type CorporateReport struct {
	CorporateGroupName string
	GroupsOfDomains    map[string][]string

	CorporateGroup *GroupData
	CommunityGroup *GroupData

	// Correlations based upon year-by-year aggregated figures for both groups
	InsertionsCorrel float64
	DeletionsCorrel  float64
	AuthorsCorrel    float64

	DomainGroupsReport           *DomainGroupsReport
	CorporateGroupSurvivalReport *GroupSurvivalReport
	CommunityGroupSurvivalReport *GroupSurvivalReport

	sqlb *db.SQLiteBackend
}

func NewCorporateReport(groupsOfDomains map[string][]string, sqlb *db.SQLiteBackend, corporateGroupName string) *CorporateReport {
	if corporateGroupName == "" {
		corporateGroupName = "Corporate"
	}

	return &CorporateReport{
		CorporateGroupName: corporateGroupName,
		GroupsOfDomains:    groupsOfDomains,
		sqlb:               sqlb,
	}
}

func (corpReport *CorporateReport) Generate() {
	domainGroupsReport := NewDomainGroupsReport(corpReport.GroupsOfDomains, corpReport.sqlb)
	domainGroupsReport.Generate()

	corpGroup := domainGroupsReport.GroupData(corpReport.CorporateGroupName)
	commGroup := domainGroupsReport.UnknownGroupData()

	insertionsCorrel, deletionsCorrel, authorsCorrel := corpGroup.Correlation(commGroup)

	corpGroupSurvival := NewGroupSurvivalReport(corpReport.sqlb, corpGroup.Authors)
	corpGroupSurvival.Generate()

	commGroupSurvival := NewGroupSurvivalReport(corpReport.sqlb, commGroup.Authors)
	commGroupSurvival.Generate()

	corpReport.DomainGroupsReport = domainGroupsReport
	corpReport.CorporateGroup = corpGroup
	corpReport.CommunityGroup = commGroup
	corpReport.InsertionsCorrel = insertionsCorrel
	corpReport.DeletionsCorrel = deletionsCorrel
	corpReport.AuthorsCorrel = authorsCorrel
	corpReport.CorporateGroupSurvivalReport = corpGroupSurvival
	corpReport.CommunityGroupSurvivalReport = commGroupSurvival
}
