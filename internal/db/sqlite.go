package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteBackend struct {
	db *sql.DB
}

func (sqlb *SQLiteBackend) DB() *sql.DB {
	return sqlb.db
}

func (sqlb *SQLiteBackend) Open(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Could not open sqlite database: %s", err)
		return err
	}

	sqlb.db = db

	return err
}

func (sqlb *SQLiteBackend) Close() error {
	err := sqlb.Close()
	if err != nil {
		log.Fatalf("Could not close sqlite database: %s", err)
	}

	return err
}

func (sqlb *SQLiteBackend) Setup() error {
	stmt := `CREATE TABLE IF NOT EXISTS commits (
			id TEXT PRIMARY KEY ON CONFLICT REPLACE,
			repo_name TEXT NOT NULL,
			author_name TEXT,
			author_email TEXT,
			author_time INT,
			committer_name TEXT,
			committer_email TEXT,
			committer_time INT,
			num_insertions INT,
			num_deletions INT,
			num_files_changed INT);
		CREATE INDEX IF NOT EXISTS index_repo_name ON commits (repo_name);
		CREATE INDEX IF NOT EXISTS index_author_name ON commits (author_name);
		CREATE INDEX IF NOT EXISTS index_author_email ON commits (author_email);
		CREATE INDEX IF NOT EXISTS index_author_time ON commits (author_time);
		CREATE INDEX IF NOT EXISTS index_committer_name ON commits (committer_name);
		CREATE INDEX IF NOT EXISTS index_committer_email ON commits (committer_email);
		CREATE INDEX IF NOT EXISTS index_committer_time ON commits (committer_time);`

	_, err := sqlb.db.Exec(stmt)
	if err != nil {
		fmt.Printf("Setup failed, received error during table creation: %s", err)
		return err
	}

	return nil
}

func (sqlb *SQLiteBackend) AddCommit(commit *common.CommitData) error {
	stmt := `INSERT INTO commits (
			id,
			repo_name,
			author_name,
			author_email,
			author_time,
			committer_name,
			committer_email,
			committer_time,
			num_insertions,
			num_deletions,
			num_files_changed
		) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11)`

	_, err := sqlb.db.Exec(stmt,
		commit.Id,
		commit.RepoName,
		commit.AuthorName,
		commit.AuthorEmail,
		commit.AuthorTime,
		commit.CommitterName,
		commit.CommitterEmail,
		commit.CommitterTime,
		commit.NumInsertions,
		commit.NumDeletions,
		commit.NumFilesChanged)

	if err != nil {
		fmt.Printf("Encountered error adding commit: %s", err)
		return err
	}

	return nil
}

func (sqlb *SQLiteBackend) Commit(commitId string) (*common.CommitData, error) {
	stmt := "SELECT * FROM commits WHERE id = ?"

	accStmt, err := sqlb.db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commit retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()

	if err != nil {
		fmt.Printf("Encountered error adding commit: %s", err)
		return nil, err
	}

	commit := new(common.CommitData)
	accStmt.QueryRow(commitId).Scan(
		&commit.Id,
		&commit.RepoName,
		&commit.AuthorName,
		&commit.AuthorEmail,
		&commit.AuthorTime,
		&commit.CommitterName,
		&commit.CommitterEmail,
		&commit.CommitterTime,
		&commit.NumInsertions,
		&commit.NumDeletions,
		&commit.NumFilesChanged,
	)

	return commit, nil
}