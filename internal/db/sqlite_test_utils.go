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

func InitTestDB() (*SQLiteBackend, error) {
	testDir, err := os.MkdirTemp("", testDirName)
	if err != nil {
		return nil, err
	}

	testDbPath := filepath.Join(testDir, testDbFileName)

	fmt.Printf("Setting up test database at %s\n", testDir)
	var sqlb = new(SQLiteBackend)
	err = sqlb.Open(testDbPath)

	if err != nil {
		log.Fatalf("Could not open database: %s", err)
		return nil, err
	}

	err = sqlb.Setup()
	if err != nil {
		log.Fatalf("Could not setup database: %s", err)
		return nil, err
	}

	return sqlb, err
}

func CleanupTestDB(sqlb *SQLiteBackend) {
	if sqlb == nil {
		return
	}

	if testDir != "" {
		os.RemoveAll(testDir)
	}
}
