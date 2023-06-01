package common

type LineChanges struct {
	NumInsertions int
	NumDeletions  int
}

type Changes struct {
	LineChanges
	NumFilesChanged int
}

type YearlyLineChangeMap map[int]LineChanges
type YearlyChangeMap map[int]Changes

// Line changes
func (lc *LineChanges) AddLineChanges(lcToAdd *LineChanges) {
	lc.NumInsertions += lcToAdd.NumInsertions
	lc.NumDeletions += lcToAdd.NumDeletions
}

func (lc *LineChanges) SubtractLineChanges(lcToSubtract *LineChanges) {
	lc.NumInsertions -= lcToSubtract.NumInsertions
	lc.NumDeletions -= lcToSubtract.NumDeletions
}

// Changes
func (changes *Changes) AddChanges(changesToAdd *Changes) {
	changes.LineChanges.AddLineChanges(&changesToAdd.LineChanges)
	changes.NumFilesChanged += changesToAdd.NumFilesChanged // FIXME: This needs to take the files into account!
}

func (changes *Changes) SubtractChanges(changesToSubtract *Changes) {
	changes.LineChanges.SubtractLineChanges(&changesToSubtract.LineChanges)
	changes.NumFilesChanged -= changesToSubtract.NumFilesChanged // FIXME: This needs to take the files into account!
}

// YearlyLineChangeMap
func (ylcm *YearlyLineChangeMap) AddLineChanges(lineChangesToAdd *LineChanges, commitYear int) {
	if changes, ok := (*ylcm)[commitYear]; ok {
		changes.AddLineChanges(lineChangesToAdd)
		(*ylcm)[commitYear] = changes
	} else {
		(*ylcm)[commitYear] = *lineChangesToAdd
	}
}

func (ylcm *YearlyLineChangeMap) SubtractLineChanges(lineChangesToSubtract *LineChanges, commitYear int) {
	if changes, ok := (*ylcm)[commitYear]; ok {
		changes.SubtractLineChanges(lineChangesToSubtract)
		(*ylcm)[commitYear] = changes
	}
}

func (ylcm *YearlyLineChangeMap) AddYearlyLineChangeMap(ylcmToAdd YearlyLineChangeMap) {
	for year, lineChangesToAdd := range ylcmToAdd {
		ylcm.AddLineChanges(&lineChangesToAdd, year)
	}
}

func (ylcm *YearlyLineChangeMap) SubtractYearlyLineChangeMap(ylcmToSubtract YearlyLineChangeMap) {
	for year, lineChangesToSubtract := range ylcmToSubtract {
		ylcm.SubtractLineChanges(&lineChangesToSubtract, year)
	}
}

// Returns insertions and deletions in two separate arraus
func (ylcm *YearlyLineChangeMap) SeparatedChangeArrays(years []int) ([]int, []int) {
	yearsToReturn := years

	if yearsToReturn == nil {
		yearsToReturn = SortedMapKeys(*ylcm)
	}

	insertionsArray := make([]int, len(yearsToReturn))
	deletionsArray := make([]int, len(yearsToReturn))

	for i, year := range yearsToReturn {
		insertionsArray[i] = (*ylcm)[year].NumInsertions
		deletionsArray[i] = (*ylcm)[year].NumDeletions
	}

	return insertionsArray, deletionsArray
}

// YearlyChangeMap
func (ycm *YearlyChangeMap) AddChanges(changesToAdd *Changes, commitYear int) {
	if changes, ok := (*ycm)[commitYear]; ok {
		changes.AddChanges(changesToAdd)
		(*ycm)[commitYear] = changes
	} else {
		(*ycm)[commitYear] = *changesToAdd
	}
}

func (ycm *YearlyChangeMap) SubtractChanges(changesToSubtract *Changes, commitYear int) {
	if changes, ok := (*ycm)[commitYear]; ok {
		changes.SubtractChanges(changesToSubtract)
		(*ycm)[commitYear] = changes
	}
}

func (ycm *YearlyChangeMap) LineChanges() YearlyLineChangeMap {
	ylcm := make(YearlyLineChangeMap, 0)

	for year, changes := range *ycm {
		ylcm[year] = changes.LineChanges
	}

	return ylcm
}
