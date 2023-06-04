package authorgroups

import (
	"testing"

	dbtesting "github.com/claucambra/commit-analysis-tool/internal/db/testing"
	"github.com/google/go-cmp/cmp"
)

func TestDomainGroupsReportGroupData(t *testing.T) {
	dbtesting.TestLogFilePath = testCommitsFile
	sqlb := dbtesting.InitTestDB(t)
	cleanup := func() { dbtesting.CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	dbtesting.IngestTestCommits(sqlb, t)

	report := NewDomainGroupsReport(testEmailGroups(), sqlb)
	report.Generate()

	testGroupData := testGroupData()
	groupData := report.GroupData(testGroupName)

	if !cmp.Equal(testGroupData, groupData) {
		t.Fatalf(`Retrieved group data does not match test group data: %s`, cmp.Diff(testGroupData, groupData))
	}

	testUnknownGroupData := testUnknownGroupData()
	unknownGroupData := report.GroupData("")

	if !cmp.Equal(testUnknownGroupData, unknownGroupData) {
		t.Fatalf(`Retrieved group data does not match test group data: %s`, cmp.Diff(testUnknownGroupData, unknownGroupData))
	}
}
