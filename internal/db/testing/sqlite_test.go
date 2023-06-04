package dbtesting

import (
	"testing"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"github.com/google/go-cmp/cmp"
)

const expectedCommitCount = 20000

func TestSqliteDbAddCommit(t *testing.T) {
	sqlb := InitTestDB(t)
	cleanup := func() { CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	parsedCommitLog := ParsedTestCommitLog(t)

	commit := parsedCommitLog[0]

	sqlb.AddCommit(commit)
	retrievedCommit, err := sqlb.Commit(commit.Id)
	if err != nil {
		t.Fatalf("Error during commit retrieval: %s", err)
	}

	if !cmp.Equal(commit, retrievedCommit) {
		t.Fatalf(`Database commit does not equal expected commit. %s`, cmp.Diff(commit, retrievedCommit))
	}
}

func TestSqliteCommits(t *testing.T) {
	sqlb := InitTestDB(t)
	cleanup := func() { CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	parsedCommits := IngestTestCommits(sqlb, t)
	receivedCommitCount := len(parsedCommits)
	if receivedCommitCount != expectedCommitCount {
		t.Fatalf("Missing commits. %+v %+v", expectedCommitCount, receivedCommitCount)
	}

	testCommits := ParsedTestCommitLog(t)

	retrievedCommits, err := sqlb.Commits()
	if err != nil {
		t.Fatalf("Could not retrieve commits in database")
	}

	CompareCommitArrays(t, testCommits, retrievedCommits)
}

func TestSqliteAuthors(t *testing.T) {
	sqlb := InitTestDB(t)
	cleanup := func() { CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	IngestTestCommits(sqlb, t)

	authors, err := sqlb.Authors()
	if err != nil {
		t.Fatalf("Received error when fetching authors: %s", err)
	}

	expectedAuthors := 247
	receivedAuthors := len(authors)

	if receivedAuthors != expectedAuthors {
		t.Fatalf("Received unexpected number of authors: expected %d, received %d", expectedAuthors, receivedAuthors)
	}
}

func TestSqliteAuthorCommits(t *testing.T) {
	sqlb := InitTestDB(t)
	cleanup := func() { CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	parsedCommits := IngestTestCommits(sqlb, t)
	receivedCommitCount := len(parsedCommits)
	if receivedCommitCount != expectedCommitCount {
		t.Fatalf("Missing commits. %+v %+v", expectedCommitCount, receivedCommitCount)
	}

	testAuthorEmail := "developer@claudiocambra.com"

	retrievedAuthorCommits, err := sqlb.AuthorCommits(testAuthorEmail)
	if err != nil {
		t.Fatalf("Could not retrieve commits for author in database")
	}

	testCommits := ParsedTestCommitLog(t)

	testAuthorCommits := make([]*common.Commit, 0)
	for _, testCommit := range testCommits {
		if testCommit.Author.Email == testAuthorEmail {
			testAuthorCommits = append(testAuthorCommits, testCommit)
		}
	}

	CompareCommitArrays(t, testAuthorCommits, retrievedAuthorCommits)
}
