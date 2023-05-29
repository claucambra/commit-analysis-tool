package common

type Changes struct {
	NumInsertions   int
	NumDeletions    int
	NumFilesChanged int
}

type YearlyChangeMap map[int]Changes
