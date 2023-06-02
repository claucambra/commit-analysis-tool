package authorgroups

import (
	"database/sql"
	"log"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

func domainChangeRows(sqlb *db.SQLiteBackend, domain string) (*sql.Rows, error) {
	stmt := "SELECT * FROM commits WHERE instr(author_email, ?) > 0"
	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commits retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()
	return accStmt.Query(domain)
}

func domainChanges(sqlb *db.SQLiteBackend, domain string) (*common.Changes, error) {
	rows, err := domainChangeRows(sqlb, domain)
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
			&commit.Author.Name,
			&commit.Author.Email,
			&commit.AuthorTime,
			&commit.Committer.Name,
			&commit.Committer.Email,
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
		LineChanges: common.LineChanges{
			NumInsertions: numInsertions,
			NumDeletions:  numDeletions,
		},
		NumFilesChanged: numFilesChanged,
	}, nil
}

func domainYearlyChanges(sqlb *db.SQLiteBackend, domain string) (common.YearlyChangeMap, error) {
	rows, err := domainChangeRows(sqlb, domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	yearBuckets := common.YearlyChangeMap{}

	for rows.Next() {
		commit := new(common.Commit)

		rows.Scan(
			&commit.Id,
			&commit.RepoName,
			&commit.Author.Name,
			&commit.Author.Email,
			&commit.AuthorTime,
			&commit.Committer.Name,
			&commit.Committer.Email,
			&commit.CommitterTime,
			&commit.NumInsertions,
			&commit.NumDeletions,
			&commit.NumFilesChanged,
		)

		commitYear := time.Unix(commit.AuthorTime, 0).Year()
		yearBuckets.AddChanges(&(commit.Changes), commitYear)
	}

	return yearBuckets, nil
}

func domainYearlyAuthors(sqlb *db.SQLiteBackend, domain string) (common.YearlyPeopleMap, error) {
	rows, err := domainChangeRows(sqlb, domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	insertedPeople := make(map[int]map[string]bool, 0)
	yearBuckets := common.YearlyPeopleMap{}

	for rows.Next() {
		commit := new(common.Commit)

		rows.Scan(
			&commit.Id,
			&commit.RepoName,
			&commit.Author.Name,
			&commit.Author.Email,
			&commit.AuthorTime,
			&commit.Committer.Name,
			&commit.Committer.Email,
			&commit.CommitterTime,
			&commit.NumInsertions,
			&commit.NumDeletions,
			&commit.NumFilesChanged,
		)

		authorEmail := commit.Author.Email
		commitYear := time.Unix(commit.AuthorTime, 0).Year()
		insertedPeopleInYear := insertedPeople[commitYear]
		personNotAlreadyAddedInYear := insertedPeopleInYear == nil || !insertedPeopleInYear[authorEmail]

		if insertedPeopleInYear == nil {
			insertedPeople[commitYear] = map[string]bool{authorEmail: true}
		} else if !insertedPeopleInYear[authorEmail] {
			insertedPeople[commitYear][authorEmail] = true
		}

		if personNotAlreadyAddedInYear {
			yearBuckets[commitYear] = append(yearBuckets[commitYear], &(commit.Author))
		}
	}

	return yearBuckets, nil
}
