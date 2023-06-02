package authorgroups

import (
	"reflect"
	"testing"

	dbtesting "github.com/claucambra/commit-analysis-tool/internal/db/testing"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

const testNumAuthors = 31
const testNumInsertions = 25205
const testNumDeletions = 17147
const testNumGroupAuthors = 5
const testGroupName = "VideoLAN"
const testGroupDomain = "videolan.org"
const testGroupInsertions = 132
const testGroupDeletions = 137
const testCommitsYear = 2023

var testCommitsFile = "../../../test/data/log.txt"

var testEmailGroups = map[string][]string{
	testGroupName: {testGroupDomain},
}

var testGroupYearlyLineChanges = map[int]*common.LineChanges{
	testCommitsYear: {
		NumInsertions: testGroupInsertions,
		NumDeletions:  testGroupDeletions,
	},
}

var testGroupYearlyAuthors = map[int][]*common.Person{
	testCommitsYear: {
		{
			Name:  "Tristan Matthews",
			Email: "tmatth@videolan.org",
		},
		{
			Name:  "Simon Latapie",
			Email: "garf@videolan.org",
		},
		{
			Name:  "Jean-Baptiste Kempf",
			Email: "jb@videolan.org",
		},
		{
			Name:  "Ilkka Ollakka",
			Email: "ileoo@videolan.org",
		},
		{
			Name:  "David Fuhrmann",
			Email: "dfuhrmann@videolan.org",
		},
	},
}

func TestNewDomainGroupsReport(t *testing.T) {
	dbtesting.TestLogFilePath = testCommitsFile
	sqlb := dbtesting.InitTestDB(t)
	cleanup := func() { dbtesting.CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	dbtesting.IngestTestCommits(sqlb, t)

	report := NewDomainGroupsReport(testEmailGroups)
	report.Generate(sqlb)

	if authorCount := report.TotalAuthors; authorCount != testNumAuthors {
		t.Fatalf("Unexpected number of authors: received %d, expected %d", authorCount, testNumAuthors)
	} else if numGroupAuthors := report.DomainTotalNumAuthors[testGroupDomain]; numGroupAuthors != testNumGroupAuthors {
		t.Fatalf("Unexpected number of domain authors: received %d, expected %d", numGroupAuthors, testNumGroupAuthors)
	}

	testGroupData := &GroupData{
		GroupName:  testGroupName,
		NumAuthors: testNumGroupAuthors,
		NumLineChanges: &common.LineChanges{
			NumInsertions: testGroupInsertions,
			NumDeletions:  testGroupDeletions,
		},
		YearlyLineChanges: testGroupYearlyLineChanges,
		YearlyAuthors:     testGroupYearlyAuthors,
		AuthorsPercent:    (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100,
		InsertionsPercent: (float32(testGroupInsertions) / float32(testNumInsertions)) * 100,
		DeletionsPercent:  (float32(testGroupDeletions) / float32(testNumDeletions)) * 100,
	}
	groupData := report.GroupData(testGroupName)

	if !reflect.DeepEqual(testGroupData, groupData) {
		t.Fatalf(`Retrieved group data does not match test group data: 
			Expected %+v
			Received %+v`, testGroupData, groupData)
	}

}
