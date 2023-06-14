package corpimpact

import (
	"strconv"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
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

func (cr *CorporateReport) Generate() {
	domainGroupsReport := authorgroups.NewDomainGroupsReport(cr.GroupsOfDomains, cr.sqlb)
	domainGroupsReport.Generate()
	cr.DomainGroupsReport = domainGroupsReport

	corpGroup := domainGroupsReport.GroupData(cr.CorporateGroupName)
	cr.CorporateGroup = corpGroup

	commGroup := domainGroupsReport.UnknownGroupData()
	cr.CommunityGroup = commGroup

	corpYearMonthInsertsMap, corpYearMonthDeletesMap, corpYearMonthAuthorsMap := cr.CorporateGroup.Commits.YearMonthCounts()
	commYearMonthInsertsMap, commYearMonthDeletesMap, commYearMonthAuthorsMap := cr.CommunityGroup.Commits.YearMonthCounts()

	cr.InsertionsCorrel = common.CorrelateYearMonthCounts(corpYearMonthInsertsMap, commYearMonthInsertsMap)
	cr.DeletionsCorrel = common.CorrelateYearMonthCounts(corpYearMonthDeletesMap, commYearMonthDeletesMap)
	cr.AuthorsCorrel = common.CorrelateYearMonthCounts(corpYearMonthAuthorsMap, commYearMonthAuthorsMap)

	corpGroupSurvival := authorgroups.NewGroupSurvivalReport(cr.sqlb, corpGroup.Authors)
	corpGroupSurvival.Generate()
	cr.CorporateGroupSurvivalReport = corpGroupSurvival

	commGroupSurvival := authorgroups.NewGroupSurvivalReport(cr.sqlb, commGroup.Authors)
	commGroupSurvival.Generate()
	cr.CommunityGroupSurvivalReport = commGroupSurvival

	corpGroupImpact := commitimpact.NewCommitImpactReport(corpGroup.Commits)
	corpGroupImpact.Generate()
	cr.CorporateCommitImpactReport = corpGroupImpact

	commGroupImpact := commitimpact.NewCommitImpactReport(commGroup.Commits)
	commGroupImpact.Generate()
	cr.CommunityCommitImpactReport = commGroupImpact
}

func (cr *CorporateReport) CSVString(name string, includeHeader bool) [][]string {
	numSurvValuesToWrite := 100
	safeCorpSurvivalValues := make([]float64, numSurvValuesToWrite)
	safeCommSurvivalValues := make([]float64, numSurvValuesToWrite)
	actualCorpSurvivalValuesLen := len(cr.CorporateGroupSurvivalReport.AuthorsSurvival)
	actualCommSurvivalValuesLen := len(cr.CommunityGroupSurvivalReport.AuthorsSurvival)

	for i := 0; i < numSurvValuesToWrite-1; i++ {
		corpSurvivalVal := 0.
		commSurvivalVal := 0.

		if i < actualCorpSurvivalValuesLen {
			corpSurvivalVal = cr.CorporateGroupSurvivalReport.AuthorsSurvival[i]
		}

		if i < actualCommSurvivalValuesLen {
			commSurvivalVal = cr.CommunityGroupSurvivalReport.AuthorsSurvival[i]
		}

		safeCorpSurvivalValues[i] = corpSurvivalVal
		safeCommSurvivalValues[i] = commSurvivalVal
	}

	csvfiedReport := []string{
		name,
		strconv.FormatInt(int64(len(cr.DomainGroupsReport.TotalCommits)), 10),
		strconv.FormatInt(int64(len(cr.DomainGroupsReport.TotalAuthors)), 10),
		strconv.FormatInt(int64(cr.DomainGroupsReport.TotalChanges.NumInsertions), 10),
		strconv.FormatInt(int64(cr.DomainGroupsReport.TotalChanges.NumDeletions), 10),
		strconv.FormatInt(int64(cr.CorporateGroup.LineChanges.NumInsertions), 10),
		strconv.FormatInt(int64(cr.CorporateGroup.LineChanges.NumDeletions), 10),
		strconv.FormatInt(int64(len(cr.CorporateGroup.Authors)), 10),
		strconv.FormatInt(int64(cr.CommunityGroup.LineChanges.NumInsertions), 10),
		strconv.FormatInt(int64(cr.CommunityGroup.LineChanges.NumDeletions), 10),
		strconv.FormatInt(int64(len(cr.CommunityGroup.Authors)), 10),
		strconv.FormatFloat(cr.CorporateGroup.InsertionsPercent, 'f', -1, 64),
		strconv.FormatFloat(cr.CorporateGroup.DeletionsPercent, 'f', -1, 64),
		strconv.FormatFloat(cr.CorporateGroup.AuthorsPercent, 'f', -1, 64),
		strconv.FormatFloat(cr.CommunityGroup.InsertionsPercent, 'f', -1, 64),
		strconv.FormatFloat(cr.CommunityGroup.DeletionsPercent, 'f', -1, 64),
		strconv.FormatFloat(cr.CommunityGroup.AuthorsPercent, 'f', -1, 64),
		strconv.FormatFloat(cr.InsertionsCorrel, 'f', -1, 64),
		strconv.FormatFloat(cr.DeletionsCorrel, 'f', -1, 64),
		strconv.FormatFloat(cr.AuthorsCorrel, 'f', -1, 64),
		strconv.FormatFloat(cr.CorporateCommitImpactReport.MeanImpact, 'f', -1, 64),
		strconv.FormatFloat(cr.CommunityCommitImpactReport.MeanImpact, 'f', -1, 64),
	}

	for i := 0; i < numSurvValuesToWrite; i++ {
		csvfiedReport = append(csvfiedReport, strconv.FormatFloat(safeCorpSurvivalValues[i], 'f', -1, 64))
	}
	for i := 0; i < numSurvValuesToWrite; i++ {
		csvfiedReport = append(csvfiedReport, strconv.FormatFloat(safeCommSurvivalValues[i], 'f', -1, 64))
	}

	var finalReport [][]string

	if includeHeader {
		header := []string{
			"name",
			"num_commits",
			"num_authors",
			"num_inserts",
			"num_deletes",
			"corp_inserts",
			"corp_deletes",
			"corp_authors",
			"comm_inserts",
			"comm_deletes",
			"comm_authors",
			"corp_insert_pc",
			"corp_delete_pc",
			"corp_authors_pc",
			"comm_insert_pc",
			"comm_delete_pc",
			"comm_authors_pc",
			"insertions_correl",
			"deletions_correl",
			"authors_correl",
			"mean_corp_impact",
			"mean_comm_impact",
		}

		for i := 0; i < numSurvValuesToWrite; i++ {
			header = append(header, "corp_surv_"+strconv.FormatInt(int64(i), 10))
		}

		for i := 0; i < numSurvValuesToWrite; i++ {
			header = append(header, "comm_surv_"+strconv.FormatInt(int64(i), 10))
		}

		finalReport = append(finalReport, header)
	}

	finalReport = append(finalReport, csvfiedReport)

	return finalReport
}

func setNumIfChildMap(childMap map[int]int, childMapIsInMap bool, childIndex int) int {
	if !childMapIsInMap {
		return 0
	} else if num, ok := childMap[childIndex]; ok {
		return num
	} else {
		return 0
	}
}

func (cr *CorporateReport) CSVChangesString(repoName string) [][]string {
	// map[Year]map[Month]NumberOfChanges
	commYearMonthInsertsMap, commYearMonthDeletesMap, commYearMonthAuthorsMap := cr.CommunityGroup.Commits.YearMonthCounts()
	corpYearMonthInsertsMap, corpYearMonthDeletesMap, corpYearMonthAuthorsMap := cr.CorporateGroup.Commits.YearMonthCounts()

	sortedCorpYears := cr.CommunityGroup.Commits.YearRange(false)
	sortedCommYears := cr.CommunityGroup.Commits.YearRange(false)
	var firstYear int

	if len(sortedCorpYears) == 0 && len(sortedCommYears) == 0 {
		return nil
	} else if len(sortedCorpYears) == 0 {
		firstYear = sortedCommYears[0]
	} else if len(sortedCommYears) == 0 {
		firstYear = sortedCorpYears[0]
	} else {
		firstYear = common.MinInt(sortedCommYears[0], sortedCorpYears[0])
	}

	maxNumYears := common.MaxInt(len(sortedCorpYears), len(sortedCommYears))

	returnArray := [][]string{
		{
			"year_month",
			"corp_insertions",
			"corp_deletions",
			"corp_authors",
			"comm_insertions",
			"comm_deletions",
			"comm_authors",
		},
	}

	for i := firstYear; i < firstYear+(maxNumYears-1); i++ {
		corpMonthInserts, yearInCorpMonthInserts := corpYearMonthInsertsMap[i]
		corpMonthDeletes, yearInCorpMonthDeletes := corpYearMonthDeletesMap[i]
		corpMonthAuthors, yearInCorpMonthAuthors := corpYearMonthAuthorsMap[i]
		commMonthInserts, yearInCommMonthInserts := commYearMonthInsertsMap[i]
		commMonthDeletes, yearInCommMonthDeletes := commYearMonthDeletesMap[i]
		commMonthAuthors, yearInCommMonthAuthors := commYearMonthAuthorsMap[i]

		for j := int(time.January); j <= int(time.December); j++ {
			corpInserts := setNumIfChildMap(corpMonthInserts, yearInCorpMonthInserts, j)
			corpDeletes := setNumIfChildMap(corpMonthDeletes, yearInCorpMonthDeletes, j)
			corpAuthors := setNumIfChildMap(corpMonthAuthors, yearInCorpMonthAuthors, j)
			commInserts := setNumIfChildMap(commMonthInserts, yearInCommMonthInserts, j)
			commDeletes := setNumIfChildMap(commMonthDeletes, yearInCommMonthDeletes, j)
			commAuthors := setNumIfChildMap(commMonthAuthors, yearInCommMonthAuthors, j)

			lineCsv := []string{
				strconv.FormatInt(int64(i), 10) + "-" + strconv.FormatInt(int64(j), 10),
				strconv.FormatInt(int64(corpInserts), 10),
				strconv.FormatInt(int64(corpDeletes), 10),
				strconv.FormatInt(int64(corpAuthors), 10),
				strconv.FormatInt(int64(commInserts), 10),
				strconv.FormatInt(int64(commDeletes), 10),
				strconv.FormatInt(int64(commAuthors), 10),
			}

			returnArray = append(returnArray, lineCsv)
		}
	}

	return returnArray
}

func (cr *CorporateReport) CSVSurvivalString(repoName string) [][]string {
	corpSurvLen := len(cr.CorporateGroupSurvivalReport.AuthorsSurvival)
	commSurvLen := len(cr.CommunityGroupSurvivalReport.AuthorsSurvival)
	totalTimeSteps := common.MaxInt(corpSurvLen, commSurvLen) + 1 // Always reach 0 for both counts

	returnArray := [][]string{
		{
			"timestep",
			"corp_survival",
			"comm_survival",
		},
	}

	for i := 0; i < totalTimeSteps; i++ {
		corpVal := 0.
		commVal := 0.

		if i < corpSurvLen {
			corpVal = cr.CorporateGroupSurvivalReport.AuthorsSurvival[i]
		}

		if i < commSurvLen {
			commVal = cr.CommunityGroupSurvivalReport.AuthorsSurvival[i]
		}

		line := []string{
			strconv.FormatInt(int64(i), 10),
			strconv.FormatFloat(corpVal, 'f', -1, 64),
			strconv.FormatFloat(commVal, 'f', -1, 64),
		}

		returnArray = append(returnArray, line)
	}

	return returnArray
}
