package dbtesting

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"github.com/claucambra/commit-analysis-tool/pkg/logread"
	"github.com/google/go-cmp/cmp"
)

const testDbFileName = "test.db"
const testDirName = "sqlite_test"

var testDir = ""
var TestLogFilePath = "../../../test/data/log.txt"

func InitTestDB(t *testing.T) *db.SQLiteBackend {
	testDir, err := os.MkdirTemp("", testDirName)
	if err != nil {
		t.Fatalf("Could not create temp test dir, received error: %s", err)
		return nil
	}

	testDbPath := filepath.Join(testDir, testDbFileName)

	log.Printf("Setting up test database at %s\n", testDir)
	var sqlb = new(db.SQLiteBackend)
	err = sqlb.Open(testDbPath)

	if err != nil {
		t.Fatalf("Could not open database: %s", err)
		return nil
	}

	err = sqlb.Setup()
	if err != nil {
		t.Fatalf("Could not setup database: %s", err)
		return nil
	}

	return sqlb
}

func ReadTestLogFile(t *testing.T) string {
	testCommitLogBytes, err := os.ReadFile(TestLogFilePath)
	if err != nil {
		t.Fatalf("Could not read test commits file")
	}

	return string(testCommitLogBytes)
}

func ParsedTestCommitLog(t *testing.T) []*common.Commit {
	testCommitLog := ReadTestLogFile(t)
	testCommits, err := logread.ParseCommitLog(testCommitLog)
	if err != nil {
		t.Fatalf("Could not parse test commit log: %s", err)
	}

	return testCommits
}

func IngestTestCommits(sqlb *db.SQLiteBackend, t *testing.T) []*common.Commit {
	parsedCommitLog := ParsedTestCommitLog(t)

	err := sqlb.AddCommits(parsedCommitLog)
	if err != nil {
		t.Fatalf("Error during test log file ingest: %s", err)
	}

	parsedCommits, err := sqlb.Commits()
	if err != nil {
		t.Fatalf("Error checking ingested commits: %s", err)
	}

	return parsedCommits
}

func CleanupTestDB(sqlb *db.SQLiteBackend) {
	if sqlb == nil {
		return
	}

	sqlb.Close()

	if testDir != "" {
		os.RemoveAll(testDir)
	}
}

func CompareCommitArrays(t *testing.T, expectedCommitArray []*common.Commit, testingCommitArray []*common.Commit) {
	numExpectedCommits := len(expectedCommitArray)
	numTestingCommits := len(testingCommitArray)

	if numExpectedCommits != numTestingCommits {
		t.Fatalf(`Expected commit count does not equal tested commit count. %s`, cmp.Diff(numExpectedCommits, numTestingCommits))
	}

	for i := 0; i < numExpectedCommits; i++ {
		expectedCommit := expectedCommitArray[i]
		testingCommit := testingCommitArray[i]

		if !cmp.Equal(expectedCommit, testingCommit) {
			t.Fatalf(`Tested commits do not equal expected commits. %s`, cmp.Diff(expectedCommit, testingCommit))
		}
	}
}
