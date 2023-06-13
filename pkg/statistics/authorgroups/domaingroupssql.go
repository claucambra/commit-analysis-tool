package authorgroups

import (
	"database/sql"
	"log"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

func domainCommitsRows(sqlb *db.SQLiteBackend, domain string) (*sql.Rows, error) {
	stmt := "SELECT * FROM commits WHERE instr(author_email, ?) > 0"
	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commits retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()
	return accStmt.Query(domain)
}

func domainCommits(sqlb *db.SQLiteBackend, domain string) ([]*common.Commit, error) {
	rows, err := domainCommitsRows(sqlb, domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	commits := []*common.Commit{}

	for rows.Next() {
		commit := sqlb.ScanRowInRowsToCommits(rows)
		commits = append(commits, commit)
	}

	return commits, nil
}

func domainLineChanges(sqlb *db.SQLiteBackend, domain string) (*common.LineChanges, error) {
	rows, err := domainCommitsRows(sqlb, domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	numInsertions := 0
	numDeletions := 0
	numFilesChanged := 0

	for rows.Next() {
		commit := sqlb.ScanRowInRowsToCommits(rows)

		numInsertions += commit.NumInsertions
		numDeletions += commit.NumDeletions
		numFilesChanged += commit.NumFilesChanged
	}

	return &common.LineChanges{
		NumInsertions: numInsertions,
		NumDeletions:  numDeletions,
	}, nil
}

func domainYearlyLineChanges(sqlb *db.SQLiteBackend, domain string) (common.YearlyLineChangeMap, error) {
	rows, err := domainCommitsRows(sqlb, domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	yearBuckets := common.YearlyLineChangeMap{}

	for rows.Next() {
		commit := sqlb.ScanRowInRowsToCommits(rows)
		commitYear := time.Unix(commit.AuthorTime, 0).UTC().Year()
		yearBuckets.AddLineChanges(&(commit.LineChanges), commitYear)
	}

	return yearBuckets, nil
}

// The years in which an author has contributed
func authorContinuousMonths(sqlb *db.SQLiteBackend, authorEmail string) (int, error) {
	authorCommits, err := sqlb.AuthorCommits(authorEmail)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return 0, err
	}

	yearsMap := map[int]map[int]bool{}
	for _, commit := range authorCommits {
		commitTime := time.Unix(commit.AuthorTime, 0).UTC()
		commitYear := commitTime.Year()
		commitMonth := int(commitTime.Month())

		if _, ok := yearsMap[commitYear]; !ok {
			yearsMap[commitYear] = map[int]bool{}
		}

		yearsMap[commitYear][commitMonth] = true
	}

	sortedYears := common.SortedMapKeys(yearsMap)
	if len(sortedYears) == 0 {
		log.Printf("Author %s active for no years, can't return number of continuous months", authorEmail)
		return 0, nil
	}

	monthCount := 0

	// Count up time, stop when found a lapse
	for i := sortedYears[0]; i <= sortedYears[len(sortedYears)-1]; i++ {
		firstMonth := int(time.January)

		if i == sortedYears[0] {
			sortedFirstYearMonths := common.SortedMapKeys(yearsMap[i])
			firstMonth = sortedFirstYearMonths[0]
		} else if _, ok := yearsMap[i]; !ok {
			return monthCount, nil
		}

		for j := firstMonth; j <= int(time.December); j++ {
			if !yearsMap[i][j] {
				return monthCount, nil
			}

			monthCount++
		}
	}

	return monthCount, nil
}
