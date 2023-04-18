package authorgroups

import (
	"testing"

	"github.com/claucambra/commit-analysis-tool/internal/db"
)

const testNumAuthors = 31
const testNumGroupAuthors = 5
const testGroupName = "VideoLAN"
const testGroupDomain = "videolan.org"

var testGroupAuthorsPercent = (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100
var testCommitsFile = "../../../test/data/log.txt"

var testEmailGroups = map[string][]string{
	testGroupName: {testGroupDomain},
}

func TestNewDomainGroupsReport(t *testing.T) {
	db.TestLogFilePath = testCommitsFile
	sqlb := db.InitTestDB(t)
	cleanup := func() { db.CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	db.IngestTestCommits(sqlb, t)

	report := NewDomainGroupsReport(testEmailGroups)
	report.Generate(sqlb)

	if authorCount := report.TotalAuthors; authorCount != testNumAuthors {
		t.Fatalf("Unexpected number of authors: received %d, expected %d", authorCount, testNumAuthors)
	} else if report.DomainNumAuthors[testGroupDomain] != testNumGroupAuthors {
		t.Fatalf("Unexpected number of domain authors: received %d, expected %d", report.DomainNumAuthors[testGroupDomain], testNumGroupAuthors)
	} else if groupPc := report.PercentageGroupAuthors(testGroupName); groupPc != testGroupAuthorsPercent {
		t.Fatalf("Unexpected percentage of group authors: received %f, expected %f", groupPc, testGroupAuthorsPercent)
	}
}
