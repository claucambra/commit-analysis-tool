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
const testOtherName = fallbackGroupName
const testNumOtherAuthors = 26
const testOtherInsertions = 25073
const testOtherDeletions = 17010
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

func testOtherAuthors() common.EmailSet {
	return common.EmailSet{
		"robux4@ycbcr.xyz":               true,
		"kerrick@kerrickstaley.com":      true,
		"khalid.masum.92@gmail.com":      true,
		"mohitmarathe23@gmail.com":       true,
		"ajanni@videolabs.io":            true,
		"andrew@crossbowffs.com":         true,
		"remi@remlab.net":                true,
		"pierre@videolabs.io":            true,
		"rom1v@videolabs.io":             true,
		"nurupo.contributions@gmail.com": true,
		"epirat07@gmail.com":             true,
		"fcvlcdev@free.fr":               true,
		"thomas@gllm.fr":                 true,
		"benjamin.arnaud@videolabs.io":   true,
		"komh@chollian.net":              true,
		"dev.asenat@posteo.net":          true,
		"linkfanel@yahoo.fr":             true,
		"guptaprince8832@gmail.com":      true,
		"johanneskauffmann@hotmail.com":  true,
		"hugo@beauzee.fr":                true,
		"loic@videolabs.io":              true,
		"git@haasn.dev":                  true,
		"fuzun54@outlook.com":            true,
		"martin@martin.st":               true,
		"umxprime@videolabs.io":          true,
		"developer@claudiocambra.com":    true,
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

func testOtherYearlyLineChanges() common.YearlyLineChangeMap {
	return common.YearlyLineChangeMap{
		2023: {
			NumInsertions: 20802,
			NumDeletions:  14618,
		},
		2022: {
			NumInsertions: 3519,
			NumDeletions:  2222,
		},
		2021: {
			NumInsertions: 254,
			NumDeletions:  103,
		},
		2020: {
			NumInsertions: 482,
			NumDeletions:  32,
		},
		2018: {
			NumInsertions: 16,
			NumDeletions:  35,
		},
	}
}

func testGroupYearlyAuthors() common.YearlyEmailMap {
	return common.YearlyEmailMap{
		testCommitsYear: testGroupAuthors(),
	}
}

func testOtherYearlyAuthors() common.YearlyEmailMap {
	return common.YearlyEmailMap{
		2023: {
			"fcvlcdev@free.fr":               true,
			"git@haasn.dev":                  true,
			"khalid.masum.92@gmail.com":      true,
			"nurupo.contributions@gmail.com": true,
			"andrew@crossbowffs.com":         true,
			"pierre@videolabs.io":            true,
			"martin@martin.st":               true,
			"kerrick@kerrickstaley.com":      true,
			"developer@claudiocambra.com":    true,
			"ajanni@videolabs.io":            true,
			"johanneskauffmann@hotmail.com":  true,
			"umxprime@videolabs.io":          true,
			"dev.asenat@posteo.net":          true,
			"mohitmarathe23@gmail.com":       true,
			"loic@videolabs.io":              true,
			"epirat07@gmail.com":             true,
			"remi@remlab.net":                true,
			"guptaprince8832@gmail.com":      true,
			"komh@chollian.net":              true,
			"thomas@gllm.fr":                 true,
			"robux4@ycbcr.xyz":               true,
			"rom1v@videolabs.io":             true,
			"fuzun54@outlook.com":            true,
			"linkfanel@yahoo.fr":             true,
			"benjamin.arnaud@videolabs.io":   true,
		},
		2022: {
			"developer@claudiocambra.com": true,
			"dev.asenat@posteo.net":       true,
			"martin@martin.st":            true,
			"hugo@beauzee.fr":             true,
			"rom1v@videolabs.io":          true,
			"robux4@ycbcr.xyz":            true,
			"ajanni@videolabs.io":         true,
			"pierre@videolabs.io":         true,
			"fcvlcdev@free.fr":            true,
		},
		2021: {
			"rom1v@videolabs.io":  true,
			"ajanni@videolabs.io": true,
		},
		2020: {
			"rom1v@videolabs.io":  true,
			"ajanni@videolabs.io": true,
		},
		2018: {
			"robux4@ycbcr.xyz": true,
		},
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

func testUnknownGroupData() *GroupData {
	return &GroupData{
		GroupName: testOtherName,
		Authors:   testOtherAuthors(),
		LineChanges: &common.LineChanges{
			NumInsertions: testOtherInsertions,
			NumDeletions:  testOtherDeletions,
		},
		YearlyLineChanges: testOtherYearlyLineChanges(),
		YearlyAuthors:     testOtherYearlyAuthors(),
		AuthorsPercent:    (float32(testNumOtherAuthors) / float32(testNumAuthors)) * 100,
		InsertionsPercent: (float32(testOtherInsertions) / float32(testNumInsertions)) * 100,
		DeletionsPercent:  (float32(testOtherDeletions) / float32(testNumDeletions)) * 100,
	}
}
