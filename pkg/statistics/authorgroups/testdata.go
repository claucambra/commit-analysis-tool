package authorgroups

import "github.com/claucambra/commit-analysis-tool/pkg/common"

const testCommitsFile = "../../../test/data/log.txt"

const testNumAuthors = 31
const testNumInsertions = 25205
const testNumDeletions = 17147
const testNumGroupAuthors = 5
const testGroupName = "VideoLAN"
const testGroupDomain = "videolan.org"
const testGroupInsertions = 132
const testGroupDeletions = 137
const testCommitsYear = 2023

var testEmailGroups = map[string][]string{
	testGroupName: {testGroupDomain},
}

var testGroupAuthors = common.EmailSet{
	"tmatth@videolan.org":    true,
	"garf@videolan.org":      true,
	"jb@videolan.org":        true,
	"ileoo@videolan.org":     true,
	"dfuhrmann@videolan.org": true,
}

var testGroupYearlyLineChanges = common.YearlyLineChangeMap{
	testCommitsYear: {
		NumInsertions: testGroupInsertions,
		NumDeletions:  testGroupDeletions,
	},
}

var testGroupYearlyAuthors = common.YearlyEmailMap{
	testCommitsYear: testGroupAuthors,
}

var testGroupData = &GroupData{
	GroupName: testGroupName,
	Authors:   testGroupAuthors,
	LineChanges: &common.LineChanges{
		NumInsertions: testGroupInsertions,
		NumDeletions:  testGroupDeletions,
	},
	YearlyLineChanges: testGroupYearlyLineChanges,
	YearlyAuthors:     testGroupYearlyAuthors,
	AuthorsPercent:    (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100,
	InsertionsPercent: (float32(testGroupInsertions) / float32(testNumInsertions)) * 100,
	DeletionsPercent:  (float32(testGroupDeletions) / float32(testNumDeletions)) * 100,
}
