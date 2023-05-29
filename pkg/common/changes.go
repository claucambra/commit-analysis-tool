package common

type LineChanges struct {
	NumInsertions int
	NumDeletions  int
}

type Changes struct {
	LineChanges

	NumFilesChanged int
}

type YearlyChangeMap map[int]Changes

func (lc *LineChanges) AddLineChanges(lcToAdd *LineChanges) {
	lc.NumInsertions += lcToAdd.NumInsertions
	lc.NumDeletions += lcToAdd.NumDeletions
}

func (lc *LineChanges) SubtractLineChanges(lcToSubtract *LineChanges) {
	lc.NumInsertions -= lcToSubtract.NumInsertions
	lc.NumDeletions -= lcToSubtract.NumDeletions
}

func (changes *Changes) AddChanges(changesToAdd *Changes) {
	changes.LineChanges.AddLineChanges(&changesToAdd.LineChanges)
	changes.NumFilesChanged += changesToAdd.NumDeletions // FIXME: This needs to take the files into account!
}

func (changes *Changes) SubtractChanges(changesToSubtract *Changes) {
	changes.LineChanges.SubtractLineChanges(&changesToSubtract.LineChanges)
	changes.NumFilesChanged += changesToSubtract.NumDeletions // FIXME: This needs to take the files into account!
}

func (ycm *YearlyChangeMap) AddChanges(changesToAdd *Changes, commitYear int) {
	if changes, ok := (*ycm)[commitYear]; ok {
		changes.AddChanges(changesToAdd)
		(*ycm)[commitYear] = changes
	} else {
		(*ycm)[commitYear] = Changes{
			LineChanges: LineChanges{
				NumInsertions: changes.NumInsertions,
				NumDeletions:  changes.NumDeletions,
			},
			NumFilesChanged: changes.NumFilesChanged,
		}
	}
}

func (ycm *YearlyChangeMap) SubtractChanges(changesToAdd *Changes, commitYear int) {
	if changes, ok := (*ycm)[commitYear]; ok {
		changes.SubtractChanges(changesToAdd)
		(*ycm)[commitYear] = changes
	}
}
