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

func TestSqliteDbAddCommit(t *testing.T) {
	sqlb := InitTestDB(t)
	cleanup := func() { CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	testCommitLogBytes, err := os.ReadFile("../../test/data/log.txt")
	if err != nil {
		t.Fatalf("Could not read test commits file")
	}

	testCommitLog := string(testCommitLogBytes)
	parsedCommitLog, err := logread.ParseCommitLog(testCommitLog)
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
