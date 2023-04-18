package db

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/claucambra/commit-analysis-tool/pkg/logread"
)

const testDbFileName = "test.db"
const testDirName = "sqlite_test"

var testDir = ""

func InitTestDB(t *testing.T) *SQLiteBackend {
	testDir, err := os.MkdirTemp("", testDirName)
	if err != nil {
		t.Fatalf("Could not create temp test dir, received error: %s", err)
		return nil
	}

	testDbPath := filepath.Join(testDir, testDbFileName)

	fmt.Printf("Setting up test database at %s\n", testDir)
	var sqlb = new(SQLiteBackend)
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

func IngestTestCommits(sqlb *SQLiteBackend, t *testing.T) {
	testCommitLogBytes, err := os.ReadFile("../../test/data/log.txt")
	if err != nil {
		t.Fatalf("Could not read test commits file")
	}

	testCommitLog := string(testCommitLogBytes)
	parsedCommitLog, err := logread.ParseCommitLog(testCommitLog)
	if err != nil {
		t.Fatalf("Error during test log file parsing: %s", err)
	}

	err = sqlb.AddCommits(parsedCommitLog)
	if err != nil {
		t.Fatalf("Error during test log file ingest: %s", err)
	}
}

func CleanupTestDB(sqlb *SQLiteBackend) {
	if sqlb == nil {
		return
	}

	sqlb.Close()

	if testDir != "" {
		os.RemoveAll(testDir)
	}
}
