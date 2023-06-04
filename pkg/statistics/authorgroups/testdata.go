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

func testEmailGroups() map[string][]string {
	return map[string][]string{
		testGroupName: {testGroupDomain},
	}
}

func testGroupAuthors() common.EmailSet {
	return common.EmailSet{
		"tmatth@videolan.org":    true,
		"garf@videolan.org":      true,
		"jb@videolan.org":        true,
		"ileoo@videolan.org":     true,
		"dfuhrmann@videolan.org": true,
	}
}

func testGroupYearlyLineChanges() common.YearlyLineChangeMap {
	return common.YearlyLineChangeMap{
		testCommitsYear: {
			NumInsertions: testGroupInsertions,
			NumDeletions:  testGroupDeletions,
		},
	}
}

func testGroupYearlyAuthors() common.YearlyEmailMap {
	return common.YearlyEmailMap{
		testCommitsYear: testGroupAuthors(),
	}
}

func testGroupData() *GroupData {
	return &GroupData{
		GroupName: testGroupName,
		Authors:   testGroupAuthors(),
		LineChanges: &common.LineChanges{
			NumInsertions: testGroupInsertions,
			NumDeletions:  testGroupDeletions,
		},
		YearlyLineChanges: testGroupYearlyLineChanges(),
		YearlyAuthors:     testGroupYearlyAuthors(),
		AuthorsPercent:    (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100,
		InsertionsPercent: (float32(testGroupInsertions) / float32(testNumInsertions)) * 100,
		DeletionsPercent:  (float32(testGroupDeletions) / float32(testNumDeletions)) * 100,
	}
}
