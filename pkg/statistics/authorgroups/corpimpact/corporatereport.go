package corpimpact

import (
	"strconv"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/authorgroups"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/commitimpact"
)

type CorporateReport struct {
	CorporateGroupName string
	GroupsOfDomains    map[string][]string

	CorporateGroup *authorgroups.GroupData
	CommunityGroup *authorgroups.GroupData

	// Correlations based upon year-by-year aggregated figures for both groups
	InsertionsCorrel float64
	DeletionsCorrel  float64
	AuthorsCorrel    float64

	DomainGroupsReport           *authorgroups.DomainGroupsReport
	CorporateGroupSurvivalReport *authorgroups.GroupSurvivalReport
	CommunityGroupSurvivalReport *authorgroups.GroupSurvivalReport

	CorporateCommitImpactReport *commitimpact.CommitImpactReport
	CommunityCommitImpactReport *commitimpact.CommitImpactReport

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
	domainGroupsReport := authorgroups.NewDomainGroupsReport(corpReport.GroupsOfDomains, corpReport.sqlb)
	domainGroupsReport.Generate()
	corpReport.DomainGroupsReport = domainGroupsReport

	corpGroup := domainGroupsReport.GroupData(corpReport.CorporateGroupName)
	corpReport.CorporateGroup = corpGroup

	commGroup := domainGroupsReport.UnknownGroupData()
	corpReport.CommunityGroup = commGroup

	insertionsCorrel, deletionsCorrel, authorsCorrel := corpGroup.Correlation(commGroup)
	corpReport.InsertionsCorrel = insertionsCorrel
	corpReport.DeletionsCorrel = deletionsCorrel
	corpReport.AuthorsCorrel = authorsCorrel

	corpGroupSurvival := authorgroups.NewGroupSurvivalReport(corpReport.sqlb, corpGroup.Authors)
	corpGroupSurvival.Generate()
	corpReport.CorporateGroupSurvivalReport = corpGroupSurvival

	commGroupSurvival := authorgroups.NewGroupSurvivalReport(corpReport.sqlb, commGroup.Authors)
	commGroupSurvival.Generate()
	corpReport.CommunityGroupSurvivalReport = commGroupSurvival

	corpGroupImpact := commitimpact.NewCommitImpactReport(corpGroup.Commits)
	corpGroupImpact.Generate()
	corpReport.CorporateCommitImpactReport = corpGroupImpact

	commGroupImpact := commitimpact.NewCommitImpactReport(commGroup.Commits)
	commGroupImpact.Generate()
	corpReport.CommunityCommitImpactReport = commGroupImpact
}

func (cr *CorporateReport) CSVString(name string, includeHeader bool) [][]string {
	numSurvValuesToWrite := 8
	safeCorpSurvivalValues := make([]float64, numSurvValuesToWrite)
	safeCommSurvivalValues := make([]float64, numSurvValuesToWrite)
	actualCorpSurvivalValuesLen := len(cr.CorporateGroupSurvivalReport.AuthorsSurvival)
	actualCommSurvivalValuesLen := len(cr.CommunityGroupSurvivalReport.AuthorsSurvival)

	for i := 0; i < numSurvValuesToWrite-1; i++ {
		corpSurvivalVal := 0.
		commSurvivalVal := 0.

		if i > actualCorpSurvivalValuesLen-1 {
			corpSurvivalVal = cr.CorporateGroupSurvivalReport.AuthorsSurvival[i]
		}

		if i > actualCommSurvivalValuesLen-1 {
			commSurvivalVal = cr.CommunityGroupSurvivalReport.AuthorsSurvival[i]
		}

		safeCorpSurvivalValues[i] = corpSurvivalVal
		safeCommSurvivalValues[i] = commSurvivalVal
	}

	csvfiedReport := []string{
		name,
		strconv.FormatInt(int64(len(cr.DomainGroupsReport.TotalCommits)), 10),
		strconv.FormatFloat(cr.InsertionsCorrel, 'E', -1, 64),
		strconv.FormatFloat(cr.DeletionsCorrel, 'E', -1, 64),
		strconv.FormatFloat(cr.AuthorsCorrel, 'E', -1, 64),
		strconv.FormatFloat(cr.CorporateCommitImpactReport.MeanImpact, 'E', -1, 64),
		strconv.FormatFloat(cr.CommunityCommitImpactReport.MeanImpact, 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[0], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[1], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[2], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[3], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[4], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[5], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[6], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[7], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[8], 'E', -1, 64),
		strconv.FormatFloat(safeCorpSurvivalValues[9], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[0], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[1], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[2], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[3], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[4], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[5], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[6], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[7], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[8], 'E', -1, 64),
		strconv.FormatFloat(safeCommSurvivalValues[9], 'E', -1, 64),
	}

	var finalReport [][]string

	if includeHeader {
		header := []string{
			"name",
			"num_commits",
			"insertions_correl",
			"deletions_correl",
			"authors_correl",
			"mean_corp_impact",
			"mean_comm_impact",
			"corp_surv_0",
			"corp_surv_1",
			"corp_surv_2",
			"corp_surv_3",
			"corp_surv_4",
			"corp_surv_5",
			"corp_surv_6",
			"corp_surv_7",
			"corp_surv_8",
			"corp_surv_9",
			"comm_surv_0",
			"comm_surv_1",
			"comm_surv_2",
			"comm_surv_3",
			"comm_surv_4",
			"comm_surv_5",
			"comm_surv_6",
			"comm_surv_7",
			"comm_surv_8",
			"comm_surv_9",
		}

		finalReport = append(finalReport, header)
	}

	finalReport = append(finalReport, csvfiedReport)

	return finalReport
}
