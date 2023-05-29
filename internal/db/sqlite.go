package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteBackend struct {
	Db *sql.DB
}

func (sqlb *SQLiteBackend) DB() *sql.DB {
	return sqlb.Db
}

func (sqlb *SQLiteBackend) Open(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Could not open sqlite database: %s", err)
		return err
	}

	sqlb.Db = db

	return err
}

func (sqlb *SQLiteBackend) Close() error {
	err := sqlb.Db.Close()
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

	_, err := sqlb.Db.Exec(stmt)
	if err != nil {
		log.Fatalf("Setup failed, received error during table creation: %s", err)
		return err
	}

	return nil
}

func (sqlb *SQLiteBackend) AddCommit(commit *common.Commit) error {
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

	_, err := sqlb.Db.Exec(stmt,
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
		log.Printf("Encountered error adding commit: %s", err)
		return err
	}

	return nil
}

func (sqlb *SQLiteBackend) AddCommits(commits []*common.Commit) error {
	for _, commit := range commits {
		err := sqlb.AddCommit(commit)

		if err != nil {
			log.Fatalf("Error adding commit: %s", err)
			return err
		}
	}

	return nil
}

func (sqlb *SQLiteBackend) Commit(commitId string) (*common.Commit, error) {
	stmt := "SELECT * FROM commits WHERE id = ?"

	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commit retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()

	commit := new(common.Commit)
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

func (sqlb *SQLiteBackend) Commits() ([]*common.Commit, error) {
	stmt := "SELECT * FROM commits"
	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commits retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()

	rows, err := accStmt.Query()
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	commits := make([]*common.Commit, 0)
	for rows.Next() {
		commit := new(common.Commit)
		rows.Scan(
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

		commits = append(commits, commit)
	}

	return commits, nil
}

func (sqlb *SQLiteBackend) Authors() ([]string, error) {
	stmt := "SELECT DISTINCT author_email FROM commits"
	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commits retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()

	rows, err := accStmt.Query()
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	authors := []string{}
	for rows.Next() {
		author := new(string)
		rows.Scan(author)
		authors = append(authors, *author)
	}

	return authors, nil
}

func (sqlb *SQLiteBackend) AuthorCommits(author string) ([]*common.Commit, error) {
	stmt := "SELECT * FROM commits WHERE author_email = ?"
	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commits retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()

	rows, err := accStmt.Query(author)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	commits := make([]*common.Commit, 0)
	for rows.Next() {
		commit := new(common.Commit)
		rows.Scan(
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

		commits = append(commits, commit)
	}

	return commits, nil
}

func (sqlb *SQLiteBackend) DomainChangeRows(domain string) (*sql.Rows, error) {
	stmt := "SELECT * FROM commits WHERE instr(author_email, ?) > 0"
	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commits retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()
	return accStmt.Query(domain)
}

func (sqlb *SQLiteBackend) DomainChanges(domain string) (*common.Changes, error) {
	rows, err := sqlb.DomainChangeRows(domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	numInsertions := 0
	numDeletions := 0
	numFilesChanged := 0

	for rows.Next() {
		commit := new(common.Commit)
		rows.Scan(
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

		numInsertions += commit.NumInsertions
		numDeletions += commit.NumDeletions
		numFilesChanged += commit.NumFilesChanged
	}

	return &common.Changes{
		NumInsertions:   numInsertions,
		NumDeletions:    numDeletions,
		NumFilesChanged: numFilesChanged,
	}, nil
}

func (sqlb *SQLiteBackend) DomainYearlyChanges(domain string) (common.YearlyChangeMap, error) {
	rows, err := sqlb.DomainChangeRows(domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	yearBuckets := make(common.YearlyChangeMap)

	for rows.Next() {
		commit := new(common.Commit)

		rows.Scan(
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

		commitYear := time.Unix(commit.AuthorTime, 0).Year()

		if _, ok := yearBuckets[commitYear]; !ok {
			yearBuckets[commitYear] = common.Changes{
				NumInsertions:   commit.NumInsertions,
				NumDeletions:    commit.NumDeletions,
				NumFilesChanged: commit.NumFilesChanged,
			}
		} else {
			existingChanges := yearBuckets[commitYear]

			existingChanges.NumInsertions += commit.NumInsertions
			existingChanges.NumDeletions += commit.NumDeletions
			existingChanges.NumFilesChanged += commit.NumFilesChanged

			yearBuckets[commitYear] = existingChanges
		}
	}

	return yearBuckets, nil
}
