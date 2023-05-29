package common

type Changes struct {
	NumInsertions   int
	NumDeletions    int
	NumFilesChanged int
}

type YearlyChangeMap map[int]Changes

func (changes *Changes) AddChanges(changesToAdd *Changes) {
	changes.NumInsertions += changesToAdd.NumInsertions
	changes.NumDeletions += changesToAdd.NumDeletions
	changes.NumFilesChanged += changesToAdd.NumDeletions // FIXME: This needs to take the files into account!
}
