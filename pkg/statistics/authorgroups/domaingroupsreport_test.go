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

	groupData := report.GroupData(testGroupName)
	expectedGroupData := testGroupData(t)

	if !cmp.Equal(expectedGroupData, groupData) {
		t.Fatalf(`Retrieved group data does not match test group data: %s`, cmp.Diff(expectedGroupData, groupData))
	}

	unknownGroupData := report.GroupData("")
	expectedUnknownGroupData := testUnknownGroupData(t)

	if !cmp.Equal(expectedUnknownGroupData, unknownGroupData) {
		t.Fatalf(`Retrieved group data does not match test group data: %s`, cmp.Diff(expectedUnknownGroupData, unknownGroupData))
	}
}
