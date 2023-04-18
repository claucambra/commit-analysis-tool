package authorgroups

import (
	"fmt"
	"testing"
	"time"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

const testNumAuthors = 3
const testGrouplessName = "unknown"
const testNumGrouplessAuthors = 1
const testNumGroupAuthors = 2
const testGroupName = "corporate"
const testGrouplessDomain = "claudiocambra.com"
const testGroupDomain = "corpdomain.com"

var testGrouplessAuthorsPercent = (float32(testNumGrouplessAuthors) / float32(testNumAuthors)) * 100
var testGroupAuthorsPercent = (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100
var testCommits = []*common.CommitData{
	{
		Id:              "8e13b181b601fed7ff4fedfd22e11821c6d621fd",
		RepoName:        "test-repo",
		AuthorName:      "Claudio Cambra",
		AuthorEmail:     fmt.Sprintf("developer@%s", testGrouplessDomain),
		AuthorTime:      time.Now().Unix(),
		CommitterName:   "Claudio Cambra",
		CommitterEmail:  fmt.Sprintf("developer@%s", testGrouplessDomain),
		CommitterTime:   time.Now().Unix(),
		NumInsertions:   2,
		NumDeletions:    0,
		NumFilesChanged: 1,
	},
	{
		Id:              "7c89d21d3bede3313d20f76b18aa82a1f6eba875",
		RepoName:        "test-repo",
		AuthorName:      "Claudio Cambra",
		AuthorEmail:     fmt.Sprintf("developer@%s", testGrouplessDomain),
		AuthorTime:      time.Now().AddDate(0, 0, -1).Unix(),
		CommitterName:   "Claudio Cambra",
		CommitterEmail:  fmt.Sprintf("developer@%s", testGrouplessDomain),
		CommitterTime:   time.Now().AddDate(0, 0, -1).Unix(),
		NumInsertions:   23,
		NumDeletions:    23,
		NumFilesChanged: 2,
	},
	{
		Id:              "c0f3fbd9a6a5acd0f0142d49fae6e4d02beb05c3",
		RepoName:        "test-repo",
		AuthorName:      "Mr. Big Wig",
		AuthorEmail:     fmt.Sprintf("bigwig@%s", testGroupDomain),
		AuthorTime:      time.Now().AddDate(0, 0, -2).Unix(),
		CommitterName:   "Claudio Cambra",
		CommitterEmail:  fmt.Sprintf("developer@%s", testGrouplessDomain),
		CommitterTime:   time.Now().AddDate(0, 0, -2).Unix(),
		NumInsertions:   197,
		NumDeletions:    10,
		NumFilesChanged: 3,
	},
	{
		Id:              "37923f8d364b9b89fd5383885dc8a220580ebda5",
		RepoName:        "test-repo",
		AuthorName:      "Dr. Big Fish",
		AuthorEmail:     fmt.Sprintf("bigfish@%s", testGroupDomain),
		AuthorTime:      time.Now().AddDate(0, 0, -3).Unix(),
		CommitterName:   "Claudio Cambra",
		CommitterEmail:  fmt.Sprintf("developer@%s", testGrouplessDomain),
		CommitterTime:   time.Now().AddDate(0, 0, -3).Unix(),
		NumInsertions:   5,
		NumDeletions:    1,
		NumFilesChanged: 1,
	},
}

var testEmailGroups = map[string]string{
	testGroupName: testGroupDomain,
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
