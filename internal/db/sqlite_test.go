package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/claucambra/commit-analysis-tool/pkg/logread"
)

const testDbFileName = "test.db"
const testDirName = "sqlite_test"

var sqlb = new(SQLiteBackend)

func setupTestDB() (string, error) {
	testDir, err := os.MkdirTemp("", testDirName)
	if err != nil {
		return "", err
	}

	testDbPath := filepath.Join(testDir, testDbFileName)

	fmt.Printf("Setting up test database at %s\n", testDir)
	err = sqlb.Open(testDbPath)
	if err != nil {
		log.Fatalf("Could not open database: %s", err)
		return "", err
	}

	err = sqlb.Setup()
	if err != nil {
		log.Fatalf("Could not setup database: %s", err)
		return "", err
	}

	return testDir, err
}

func cleanup(path string) {
	os.RemoveAll(path)
}

func TestSqliteDbSetup(t *testing.T) {
	testPath, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup database: %s", err)
	}

	cleanup := func() { cleanup(testPath) }
	t.Cleanup(cleanup)

}

func TestSqliteDbAddCommit(t *testing.T) {
	testPath, _ := setupTestDB()
	cleanup := func() { cleanup(testPath) }
	t.Cleanup(cleanup)

	testCommitLogBytes, err := os.ReadFile("../../test/data/log.txt")
	if err != nil {
		t.Fatalf("Could not read test commits file")
	}

	testCommitLog := string(testCommitLogBytes)
	parsedCommitLog, err := logread.ParseCommitLog(testCommitLog, nil)
	if err != nil {
		t.Fatalf("Error during test log file parsing: %s", err)
	}

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
