package authorgroups

import (
	"reflect"
	"strings"
	"testing"

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
