package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func CleanupTestDB(sqlb *SQLiteBackend) {
	if sqlb == nil {
		return
	}

	if testDir != "" {
		os.RemoveAll(testDir)
	}
}
