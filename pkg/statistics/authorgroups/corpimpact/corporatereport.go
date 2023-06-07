package corpimpact

import (
	"strconv"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/authorgroups"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/commitimpact"
)

type YearMonthCount map[int]map[int]int

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

	// TODO: Maybe use year-month count flattened maps here too? More data
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
	numSurvValuesToWrite := 10
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

// FIXME: Just create an additive func for this
func addValInYearMonthCountMap(inMap YearMonthCount, year int, month int, value int) {
	if _, ok := inMap[year]; ok {
		common.AdditiveValueMapInsert[int, int, map[int]int](inMap[year], month, func(a int, b int) int {
			return a + b // FIXME: And here too
		}, value)
	} else {
		inMap[year] = map[int]int{month: value}
	}
}

func yearMonthlyDataMap(group *authorgroups.GroupData) (YearMonthCount, YearMonthCount, YearMonthCount) {
	yearMonthInsertsMap := YearMonthCount{}
	yearMonthDeletesMap := YearMonthCount{}
	yearMonthAuthorsMap := YearMonthCount{}

	checkAuthorInYearMonthMap := map[int]map[int]map[string]bool{}

	for _, commit := range group.Commits {
		commitTime := time.Unix(commit.AuthorTime, 0).UTC()
		commitYear := commitTime.Year()
		commitMonth := int(commitTime.Month())
		commitAuthor := commit.Author.Email

		addValInYearMonthCountMap(yearMonthInsertsMap, commitYear, commitMonth, commit.LineChanges.NumInsertions)
		addValInYearMonthCountMap(yearMonthDeletesMap, commitYear, commitMonth, commit.LineChanges.NumDeletions)

		// Make sure to only add author if not added already
		addAuthor := false

		if monthMap, ok := checkAuthorInYearMonthMap[commitYear]; !ok {
			checkAuthorInYearMonthMap[commitYear] = map[int]map[string]bool{commitMonth: {commitAuthor: true}}
			addAuthor = true
		} else if monthAuthors, ok := monthMap[commitMonth]; !ok {
			checkAuthorInYearMonthMap[commitYear][commitMonth] = map[string]bool{commitAuthor: true}
			addAuthor = true
		} else if !monthAuthors[commitAuthor] {
			checkAuthorInYearMonthMap[commitYear][commitMonth][commitAuthor] = true
			addAuthor = true
		}

		if addAuthor {
			addValInYearMonthCountMap(yearMonthAuthorsMap, commitYear, commitMonth, 1)
		}
	}

	return yearMonthInsertsMap, yearMonthDeletesMap, yearMonthAuthorsMap
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

// TODO: Move processing to DomainGroupsReport and base on SQL
func (cr *CorporateReport) CSVChangesString(repoName string) [][]string {
	// map[Year]map[Month]NumberOfChanges
	commYearMonthInsertsMap, commYearMonthDeletesMap, commYearMonthAuthorsMap := yearMonthlyDataMap(cr.CommunityGroup)
	corpYearMonthInsertsMap, corpYearMonthDeletesMap, corpYearMonthAuthorsMap := yearMonthlyDataMap(cr.CorporateGroup)

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

	sortedCorpYears := common.SortedMapKeys(cr.CorporateGroup.YearlyLineChanges)
	sortedCommYears := common.SortedMapKeys(cr.CommunityGroup.YearlyLineChanges)
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
