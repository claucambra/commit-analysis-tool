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
	grouplessGroup := report.DomainGroups[testGrouplessName]
	if group == nil || grouplessGroup == nil {
		t.Fatalf("Fetched author group was nil")
	}

	if report.AuthorCount != testNumAuthors {
		t.Fatalf("Unexpected number of authors: received %d, expected %d", report.AuthorCount, testNumAuthors)
	} else if group.AuthorCount != testNumGroupAuthors {
		t.Fatalf("Unexpected number of group authors: received %d, expected %d", group.AuthorCount, testNumGroupAuthors)
	} else if grouplessGroup.AuthorCount != testNumGrouplessAuthors {
		t.Fatalf("Unexpected number of unknown group authors: received %d, expected %d", grouplessGroup.AuthorCount, testNumGrouplessAuthors)
	} else if grouplessAuthors := report.AuthorCount - group.AuthorCount; grouplessAuthors != testNumGrouplessAuthors {
		t.Fatalf("Unexpected number of groupless authors: received %d, expected %d", grouplessAuthors, testNumGrouplessAuthors)
	} else if groupPercentage := report.GroupPercentageOfTotal(testGroupName); groupPercentage != testGroupAuthorsPercent {
		t.Fatalf("Unexpected group author percent: received %f, expected %f", groupPercentage, testGroupAuthorsPercent)
	} else if grouplessGroupPercentage := report.GroupPercentageOfTotal(testGrouplessName); grouplessGroupPercentage != testGrouplessAuthorsPercent {
		t.Fatalf("Unexpected groupless author percent: received %f, expected %f", grouplessGroupPercentage, testGrouplessAuthorsPercent)
	}
}

func TestDomainGroupsString(t *testing.T) {
	testString := "Author domain groups report\n"
	testString += "Total repository authors: 31\n"
	testString += "Number of authors by group:\n"
	testString += "\t\"unknown\":\t26 (83.870964%)\n"
	testString += "\t\tvideolabs.io:\t6\n"
	testString += "\t\tgmail.com:\t5\n"
	testString += "\t\tbeauzee.fr:\t1\n"
	testString += "\t\tchollian.net:\t1\n"
	testString += "\t\tclaudiocambra.com:\t1\n"
	testString += "\t\tcrossbowffs.com:\t1\n"
	testString += "\t\tfree.fr:\t1\n"
	testString += "\t\tgllm.fr:\t1\n"
	testString += "\t\thaasn.dev:\t1\n"
	testString += "\t\thotmail.com:\t1\n"
	testString += "\t\tkerrickstaley.com:\t1\n"
	testString += "\t\tmartin.st:\t1\n"
	testString += "\t\toutlook.com:\t1\n"
	testString += "\t\tposteo.net:\t1\n"
	testString += "\t\tremlab.net:\t1\n"
	testString += "\t\tyahoo.fr:\t1\n"
	testString += "\t\tycbcr.xyz:\t1\n"
	testString += "\t\"VideoLAN\":\t5 (16.129032%)\n"
	testString += "\t\tvideolan.org:\t5\n"

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
