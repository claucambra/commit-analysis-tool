package statistics

import (
	"testing"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/git"
)

var testCommits = []*git.CommitData{
	{
		Id:              "8e13b181b601fed7ff4fedfd22e11821c6d621fd",
		RepoName:        "test-repo",
		AuthorName:      "Claudio Cambra",
		AuthorEmail:     "developer@claudiocambra.com",
		AuthorTime:      time.Now().Unix(),
		CommitterName:   "Claudio Cambra",
		CommitterEmail:  "developer@claudiocambra.com",
		CommitterTime:   time.Now().Unix(),
		NumInsertions:   2,
		NumDeletions:    0,
		NumFilesChanged: 1,
	},
	{
		Id:              "7c89d21d3bede3313d20f76b18aa82a1f6eba875",
		RepoName:        "test-repo",
		AuthorName:      "Claudio Cambra",
		AuthorEmail:     "developer@claudiocambra.com",
		AuthorTime:      time.Now().AddDate(0, 0, -1).Unix(),
		CommitterName:   "Claudio Cambra",
		CommitterEmail:  "developer@claudiocambra.com",
		CommitterTime:   time.Now().AddDate(0, 0, -1).Unix(),
		NumInsertions:   23,
		NumDeletions:    23,
		NumFilesChanged: 2,
	},
	{
		Id:              "c0f3fbd9a6a5acd0f0142d49fae6e4d02beb05c3",
		RepoName:        "test-repo",
		AuthorName:      "Mr. Big Wig",
		AuthorEmail:     "bigwig@corpmail.com",
		AuthorTime:      time.Now().AddDate(0, 0, -2).Unix(),
		CommitterName:   "Claudio Cambra",
		CommitterEmail:  "developer@claudiocambra.com",
		CommitterTime:   time.Now().AddDate(0, 0, -2).Unix(),
		NumInsertions:   197,
		NumDeletions:    10,
		NumFilesChanged: 3,
	},
	{
		Id:              "37923f8d364b9b89fd5383885dc8a220580ebda5",
		RepoName:        "test-repo",
		AuthorName:      "Dr. Big Fish",
		AuthorEmail:     "bigfish@corpmail.com",
		AuthorTime:      time.Now().AddDate(0, 0, -3).Unix(),
		CommitterName:   "Claudio Cambra",
		CommitterEmail:  "developer@claudiocambra.com",
		CommitterTime:   time.Now().AddDate(0, 0, -3).Unix(),
		NumInsertions:   5,
		NumDeletions:    1,
		NumFilesChanged: 1,
	},
}

func TestNewCorporateAuthorsReport(t *testing.T) {
	testCorpEmails := make(map[string]bool)
	testCorpEmails["corpmail.com"] = true

	testNumAuthors := 3
	testNumCorpAuthors := 2
	testCorpAuthorsPercent := (float32(testNumCorpAuthors) / float32(testNumAuthors)) * 100

	report := NewCorpAuthorsReport(testCommits, testCorpEmails)
	if report.TotalAuthors != testNumAuthors {
		t.Fatalf("Unexpected number of authors: received %d, expected %d", report.TotalAuthors, testNumAuthors)
	} else if report.NumCorpAuthors != testNumCorpAuthors {
		t.Fatalf("Unexpected number of corporate authors: received %d, expected %d", report.NumCorpAuthors, testNumCorpAuthors)
	} else if report.CorpAuthorsPercent != testCorpAuthorsPercent {
		t.Fatalf("Unexpected corporate author percent: received %f, expected %f", report.CorpAuthorsPercent, testCorpAuthorsPercent)
	}
}
