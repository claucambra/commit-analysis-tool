package authorgroups

import (
	"github.com/claucambra/commit-analysis-tool/internal/db"
	"gonum.org/v1/gonum/stat"
)

type CorporateReport struct {
	CorporateGroup *GroupData
	CommunityGroup *GroupData

	InsertionsCorrel float64
	DeletionsCorrel  float64

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

	yearsInBothGroups := []int{}
	for year := range corpGroup.YearlyLineChanges {
		_, ok := commGroup.YearlyLineChanges[year]
		if ok {
			yearsInBothGroups = append(yearsInBothGroups, year)
		}
	}

	corpYearlyInsertionsArr, corpYearlyDeletionsArr := corpGroup.YearlyLineChanges.SeparatedChangeArrays(yearsInBothGroups)
	commYearlyInsertionsArr, commYearlyDeletionsArr := commGroup.YearlyLineChanges.SeparatedChangeArrays(yearsInBothGroups)

	floatCorpInserts := make([]float64, len(yearsInBothGroups))
	floatCorpDeletes := make([]float64, len(yearsInBothGroups))
	floatCommInserts := make([]float64, len(yearsInBothGroups))
	floatCommDeletes := make([]float64, len(yearsInBothGroups))

	for i := 0; i < len(corpGroup.YearlyLineChanges); i++ {
		floatCorpInserts[i] = float64(corpYearlyInsertionsArr[i])
		floatCorpDeletes[i] = float64(corpYearlyDeletionsArr[i])
	}
	for i := 0; i < len(corpGroup.YearlyLineChanges); i++ {
		floatCommInserts[i] = float64(commYearlyInsertionsArr[i])
		floatCommDeletes[i] = float64(commYearlyDeletionsArr[i])
	}

	insertionsCorrel := stat.Correlation(floatCorpInserts, floatCommInserts, nil)
	deletionsCorrel := stat.Correlation(floatCorpDeletes, floatCommDeletes, nil)

	return &CorporateReport{
		CorporateGroup: corpGroup,
		CommunityGroup: commGroup,

		InsertionsCorrel: insertionsCorrel,
		DeletionsCorrel:  deletionsCorrel,

		domainGroupsReport: domainGroupsReport,
	}
}
