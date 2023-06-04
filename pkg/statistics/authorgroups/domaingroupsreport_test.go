package authorgroups

import (
	"reflect"
	"testing"

	dbtesting "github.com/claucambra/commit-analysis-tool/internal/db/testing"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

func TestDomainGroupsReportGroupData(t *testing.T) {
	dbtesting.TestLogFilePath = testCommitsFile
	sqlb := dbtesting.InitTestDB(t)
	cleanup := func() { dbtesting.CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	dbtesting.IngestTestCommits(sqlb, t)

	report := NewDomainGroupsReport(testEmailGroups, sqlb)
	report.Generate()

	testGroupData := &GroupData{
		GroupName: testGroupName,
		Authors:   testGroupAuthors,
		LineChanges: &common.LineChanges{
			NumInsertions: testGroupInsertions,
			NumDeletions:  testGroupDeletions,
		},
		YearlyLineChanges: testGroupYearlyLineChanges,
		YearlyAuthors:     testGroupYearlyAuthors,
		AuthorsPercent:    (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100,
		InsertionsPercent: (float32(testGroupInsertions) / float32(testNumInsertions)) * 100,
		DeletionsPercent:  (float32(testGroupDeletions) / float32(testNumDeletions)) * 100,
	}
	groupData := report.GroupData(testGroupName)

	if !cmp.Equal(testGroupData, groupData) {
		t.Fatalf(`Retrieved group data does not match test group data: %s`, cmp.Diff(testGroupData, groupData))
	}
}
