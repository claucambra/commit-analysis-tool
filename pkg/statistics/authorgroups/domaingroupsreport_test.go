package authorgroups

import (
	"os"
	"testing"

	"github.com/claucambra/commit-analysis-tool/pkg/logread"
)

const testNumAuthors = 31
const testGrouplessName = "unknown"
const testNumGrouplessAuthors = 26
const testNumGroupAuthors = 5
const testGroupName = "VideoLAN"
const testGrouplessDomain = "claudiocambra.com"
const testGroupDomain = "videolan.org"

var testGrouplessAuthorsPercent = (float32(testNumGrouplessAuthors) / float32(testNumAuthors)) * 100
var testGroupAuthorsPercent = (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100
var testCommitsFile = "../../../test/data/log.txt"
var testCommitsBytes, readFileErr = os.ReadFile(testCommitsFile)
var testCommitsString = string(testCommitsBytes)
var testCommits, readCommitsErr = logread.ParseCommitLog(testCommitsString, nil)

var testEmailGroups = map[string]string{
	testGroupName: testGroupDomain,
}

func TestCommitsFile(t *testing.T) {
	if readFileErr != nil {
		t.Fatalf("Received error on test data read: %s", readFileErr)
	}

	if readCommitsErr != nil {
		t.Fatalf("Received error on reading commits: %s", readCommitsErr)
	}
}

func TestNewDomainGroupsReport(t *testing.T) {
	report := NewDomainGroupsReport(testEmailGroups)
	report.ParseCommits(testCommits)

	group := report.DomainGroups[testGroupName]
	if group == nil {
		t.Fatalf("Fetched author group was nil")
	}

	if report.AuthorCount != testNumAuthors {
		t.Fatalf("Unexpected number of authors: received %d, expected %d", report.AuthorCount, testNumAuthors)
	} else if group.AuthorCount != testNumGroupAuthors {
		t.Fatalf("Unexpected number of group authors: received %d, expected %d", group.AuthorCount, testNumGroupAuthors)
	} else if grouplessAuthors := report.AuthorCount - group.AuthorCount; grouplessAuthors != testNumGrouplessAuthors {
		t.Fatalf("Unexpected number of groupless authors: received %d, expected %d", grouplessAuthors, testNumGrouplessAuthors)
	} else if groupPercentage := report.GroupPercentageOfTotal(testGroupName); groupPercentage != testGroupAuthorsPercent {
		t.Fatalf("Unexpected group author percent: received %f, expected %f", groupPercentage, testGroupAuthorsPercent)
	}
}

func TestDomainGroupsString(t *testing.T) {
	testString := "Author domain groups report\n"
	testString += fmt.Sprintf("Total repository authors: %d\n", testNumAuthors)

	testString += "Number of authors by group:\n"
	testString += fmt.Sprintf("\t\"%s\":\t%d (%f%%)\n", testGrouplessName, testNumGrouplessAuthors, testGrouplessAuthorsPercent)
	testString += fmt.Sprintf("\t\t%s:\t%d\n", testGrouplessDomain, testNumGrouplessAuthors)
	testString += fmt.Sprintf("\t\"%s\":\t%d (%f%%)\n", testGroupName, testNumGroupAuthors, testGroupAuthorsPercent)
	testString += fmt.Sprintf("\t\t%s:\t%d\n", testGroupDomain, testNumGroupAuthors)

	report := NewDomainGroupsReport(testEmailGroups)
	report.ParseCommits(testCommits)

	reportString := report.String()
	if reportString != testString {
		t.Fatalf(`Received stringification does not match expected.
			Received: %s
			Expected: %s`, reportString, testString)
	}

	t.Logf("Received correct stringification: %s", reportString)
}
