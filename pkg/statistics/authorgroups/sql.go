package authorgroups

import (
	"database/sql"
	"log"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

func DomainChangeRows(sqlb *db.SQLiteBackend, domain string) (*sql.Rows, error) {
	stmt := "SELECT * FROM commits WHERE instr(author_email, ?) > 0"
	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commits retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()
	return accStmt.Query(domain)
}

func DomainChanges(sqlb *db.SQLiteBackend, domain string) (*common.Changes, error) {
	rows, err := DomainChangeRows(sqlb, domain)
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

func DomainYearlyChanges(sqlb *db.SQLiteBackend, domain string) (common.YearlyChangeMap, error) {
	rows, err := DomainChangeRows(sqlb, domain)
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
