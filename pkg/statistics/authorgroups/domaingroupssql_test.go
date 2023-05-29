package authorgroups

import (
	"reflect"
	"strings"
	"testing"
	"time"

	dbtesting "github.com/claucambra/commit-analysis-tool/internal/db/testing"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
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

	retrievedDomainChanges, err := domainChanges(sqlb, testDomain)
	if err != nil {
		t.Fatalf("Error retrieving domain's changes from database")
	}

	parsedCommitLog := dbtesting.ParsedTestCommitLog(t)
	testDomainChanges := &common.Changes{
		NumInsertions:   0,
		NumDeletions:    0,
		NumFilesChanged: 0,
	}

	for _, commit := range parsedCommitLog {
		if !emailHasDomain(commit.AuthorEmail, testDomain) {
			continue
		}

		testDomainChanges.NumInsertions += commit.NumInsertions
		testDomainChanges.NumDeletions += commit.NumDeletions
		testDomainChanges.NumFilesChanged += commit.NumFilesChanged
	}

	if !reflect.DeepEqual(retrievedDomainChanges, testDomainChanges) {
		t.Fatalf(`Database domain changes do not equal expected domain changes.
			Expected: %+v
			Received: %+v`, testDomainChanges, retrievedDomainChanges)
	}
}

func TestDomainYearlyChanges(t *testing.T) {
	sqlb := dbtesting.InitTestDB(t)
	cleanup := func() { dbtesting.CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	dbtesting.IngestTestCommits(sqlb, t)

	retrievedDomainYearlyChanges, err := domainYearlyChanges(sqlb, testDomain)
	if err != nil {
		t.Fatalf("Error retrieving domain's yearly changes from database")
	}

	parsedCommitLog := dbtesting.ParsedTestCommitLog(t)
	testDomainYearlyChanges := make(common.YearlyChangeMap, 0)

	for _, commit := range parsedCommitLog {
		if !emailHasDomain(commit.AuthorEmail, testDomain) {
			continue
		}

		commitYear := time.Unix(commit.AuthorTime, 0).Year()

		if changes, ok := testDomainYearlyChanges[commitYear]; ok {
			changes.AddChanges(&(commit.Changes))
			testDomainYearlyChanges[commitYear] = changes
		} else {
			testDomainYearlyChanges[commitYear] = common.Changes{
				NumInsertions:   commit.NumInsertions,
				NumDeletions:    commit.NumDeletions,
				NumFilesChanged: commit.NumFilesChanged,
			}
		}
	}

	numTestYears := len(testDomainYearlyChanges)
	numRetrievedYears := len(retrievedDomainYearlyChanges)

	if numRetrievedYears != numTestYears {
		t.Fatalf(`Number of retrieved domain change years do not equal expected domain change years.
			Expected: %+v
			Received: %+v`, numTestYears, numRetrievedYears)
	}

	for year, testChanges := range testDomainYearlyChanges {
		retrievedChanges, ok := retrievedDomainYearlyChanges[year]

		if !ok {
			t.Fatalf(`Retrieved yearly changes does not contain the year %+v`, year)
		}

		if !reflect.DeepEqual(retrievedChanges, testChanges) {
			t.Fatalf(`Database domain changes do not equal expected domain changes.
				Expected: %+v
				Received: %+v
				Error occurred when testing results for year %+v`, testChanges, retrievedChanges, year)
		}
	}
}
