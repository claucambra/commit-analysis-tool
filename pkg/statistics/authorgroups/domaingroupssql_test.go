package authorgroups

import (
	"strings"
	"testing"
	"time"

	dbtesting "github.com/claucambra/commit-analysis-tool/internal/db/testing"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"github.com/google/go-cmp/cmp"
)

var testDomain = "claudiocambra.com"

func emailHasDomain(email string, domain string) bool {
	splitEmail := strings.Split(email, "@")

	if len(splitEmail) != 2 {
		return false
	}

	domainFromEmail := splitEmail[1]
	return domainFromEmail == domain
}

func TestDomainChanges(t *testing.T) {
	sqlb := dbtesting.InitTestDB(t)
	cleanup := func() { dbtesting.CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	dbtesting.IngestTestCommits(sqlb, t)

	retrievedDomainChanges, err := domainLineChanges(sqlb, testDomain)
	if err != nil {
		t.Fatalf("Error retrieving domain's changes from database")
	}

	parsedCommitLog := dbtesting.ParsedTestCommitLog(t)
	testDomainLineChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}

	for _, commit := range parsedCommitLog {
		if !emailHasDomain(commit.Author.Email, testDomain) {
			continue
		}

		testDomainLineChanges.NumInsertions += commit.NumInsertions
		testDomainLineChanges.NumDeletions += commit.NumDeletions
	}

	if !cmp.Equal(retrievedDomainChanges, testDomainLineChanges) {
		t.Fatalf(`Database domain changes do not equal expected domain changes. %s`, cmp.Diff(testDomainLineChanges, retrievedDomainChanges))
	}
}

func TestDomainYearlyChanges(t *testing.T) {
	sqlb := dbtesting.InitTestDB(t)
	cleanup := func() { dbtesting.CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	dbtesting.IngestTestCommits(sqlb, t)

	retrievedDomainYearlyLineChanges, err := domainYearlyLineChanges(sqlb, testDomain)
	if err != nil {
		t.Fatalf("Error retrieving domain's yearly changes from database")
	}

	parsedCommitLog := dbtesting.ParsedTestCommitLog(t)
	testDomainYearlyLineChanges := make(common.YearlyLineChangeMap, 0)

	for _, commit := range parsedCommitLog {
		if !emailHasDomain(commit.Author.Email, testDomain) {
			continue
		}

		commitYear := time.Unix(commit.AuthorTime, 0).Year()
		testDomainYearlyLineChanges.AddLineChanges(&(commit.LineChanges), commitYear)
	}

	numTestYears := len(testDomainYearlyLineChanges)
	numRetrievedYears := len(retrievedDomainYearlyLineChanges)

	if numRetrievedYears != numTestYears {
		t.Fatalf(`Number of retrieved domain change years do not equal expected domain change years.
			Expected: %+v
			Received: %+v`, numTestYears, numRetrievedYears)
	}

	for year, testLineChanges := range testDomainYearlyLineChanges {
		retrievedChanges, ok := retrievedDomainYearlyLineChanges[year]

		if !ok {
			t.Fatalf(`Retrieved yearly changes does not contain the year %+v`, year)
		}

		if !cmp.Equal(retrievedChanges, testLineChanges) {
			t.Fatalf(`Database domain changes do not equal expected domain changes.
				%s
				Error occurred when testing results for year %+v`, cmp.Diff(testLineChanges, retrievedChanges), year)
		}
	}
}
