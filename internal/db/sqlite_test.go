package db

import (
	"reflect"
	"testing"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

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

	if !reflect.DeepEqual(commit, retrievedCommit) {
		t.Fatalf(`Database commit does not equal expected commit.
			Expected: %+v
			Received: %+v`, commit, retrievedCommit)
	}
}

func TestSqliteCommits(t *testing.T) {
	sqlb := InitTestDB(t)
	cleanup := func() { CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	IngestTestCommits(sqlb, t)

	testCommits := ParsedTestCommitLog(t)

	retrievedCommits, err := sqlb.Commits()
	if err != nil {
		t.Fatalf("Could not retrieve commits in database")
	}

	numTestCommits := len(testCommits)
	numRetrievedCommits := len(retrievedCommits)

	if numRetrievedCommits != numTestCommits {
		t.Fatalf(`Database commit count does not equal expected commit count.
			Expected: %+v commits
			Received: %+v commits`, numTestCommits, numRetrievedCommits)
	}

	for i := 0; i < numTestCommits; i++ {
		testCommit := testCommits[i]
		retrievedCommit := retrievedCommits[i]

		if !reflect.DeepEqual(testCommit, retrievedCommit) {
			t.Fatalf(`Database commits does not equal expected commits.
				Expected: %+v
				Received: %+v`, testCommit, retrievedCommit)
		}
	}
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

	if len(authors) != 31 {
		t.Fatalf("Received unexpected number of authors: expected %d, received %d", 1, len(authors))
	}
}

func TestSqliteAuthorCommits(t *testing.T) {
	sqlb := InitTestDB(t)
	cleanup := func() { CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	IngestTestCommits(sqlb, t)

	testAuthorEmail := "developer@claudiocambra.com"

	retrievedAuthorCommits, err := sqlb.AuthorCommits(testAuthorEmail)
	if err != nil {
		t.Fatalf("Could not retrieve commits for author in database")
	}

	testCommits := ParsedTestCommitLog(t)

	testAuthorCommits := make([]*common.Commit, 0)
	for _, testCommit := range testCommits {
		if testCommit.AuthorEmail == testAuthorEmail {
			testAuthorCommits = append(testAuthorCommits, testCommit)
		}
	}

	numTestAuthorCommits := len(testAuthorCommits)
	numRetrievedAuthorCommits := len(retrievedAuthorCommits)

	if numRetrievedAuthorCommits != numTestAuthorCommits {
		t.Fatalf(`Database commit count does not equal expected commit count.
			Expected: %+v commits
			Received: %+v commits`, numTestAuthorCommits, numRetrievedAuthorCommits)
	}

	for i := 0; i < numTestAuthorCommits; i++ {
		testAuthorCommit := testAuthorCommits[i]
		retrievedAuthorCommit := retrievedAuthorCommits[i]

		if !reflect.DeepEqual(testAuthorCommit, retrievedAuthorCommit) {
			t.Fatalf(`Database commits does not equal expected commits.
				Expected: %+v
				Received: %+v`, testAuthorCommit, retrievedAuthorCommit)
		}
	}
}
