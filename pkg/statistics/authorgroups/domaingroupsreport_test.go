package authorgroups

import (
	"os"
	"testing"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/logread"
)

const testNumAuthors = 31
const testGrouplessName = "unknown"
const testNumGrouplessAuthors = 26
const testNumGroupAuthors = 5
const testGroupName = "VideoLAN"
const testGroupDomain = "videolan.org"

var testGrouplessAuthorsPercent = (float32(testNumGrouplessAuthors) / float32(testNumAuthors)) * 100
var testGroupAuthorsPercent = (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100
var testCommitsFile = "../../../test/data/log.txt"
var testCommitsBytes, readFileErr = os.ReadFile(testCommitsFile)
var testCommitsString = string(testCommitsBytes)
var testCommits, readCommitsErr = logread.ParseCommitLog(testCommitsString)

var testEmailGroups = map[string][]string{
	testGroupName: {testGroupDomain},
}

func TestCommitsFile(t *testing.T) {
	db.TestLogFilePath = "../../../test/data/log.txt"
	if readFileErr != nil {
		t.Fatalf("Received error on test data read: %s", readFileErr)
	}

	if readCommitsErr != nil {
		t.Fatalf("Received error on reading commits: %s", readCommitsErr)
	}
}

func TestNewDomainGroupsReport(t *testing.T) {
	db.TestLogFilePath = "../../../test/data/log.txt"
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
