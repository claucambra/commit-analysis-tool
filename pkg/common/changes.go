package common

type LineChanges struct {
	NumInsertions int
	NumDeletions  int
}

type Changes struct {
	LineChanges
	NumFilesChanged int
}

type YearlyLineChangeMap map[int]*LineChanges
type YearlyChangeMap map[int]*Changes

func AddLineChanges(a *LineChanges, b *LineChanges) *LineChanges {
	return &LineChanges{
		NumInsertions: a.NumInsertions + b.NumInsertions,
		NumDeletions:  a.NumDeletions + b.NumDeletions,
	}
}

func SubtractLineChanges(a *LineChanges, b *LineChanges) (*LineChanges, bool) {
	subtractedInsertions := MaxInt(a.NumInsertions-b.NumInsertions, 0)
	subtractedDeletions := MaxInt(a.NumDeletions-b.NumDeletions, 0)

	lineChanges := &LineChanges{
		NumInsertions: subtractedInsertions,
		NumDeletions:  subtractedDeletions,
	}

	emptyOrInvalid := subtractedInsertions <= 0 && subtractedDeletions <= 0
	return lineChanges, emptyOrInvalid
}

func AddChanges(a *Changes, b *Changes) *Changes {
	return &Changes{
		LineChanges:     *AddLineChanges(&a.LineChanges, &b.LineChanges),
		NumFilesChanged: a.NumFilesChanged + b.NumFilesChanged,
	}
}

func SubtractChanges(a *Changes, b *Changes) (*Changes, bool) {
	subtractedLineChanges, lcEmptyOrInvalid := SubtractLineChanges(&a.LineChanges, &b.LineChanges)
	subtractedFilesChanged := MaxInt(a.NumFilesChanged-b.NumFilesChanged, 0)
	changes := &Changes{
		LineChanges:     *subtractedLineChanges,
		NumFilesChanged: subtractedFilesChanged,
	}

	emptyOrInvalid := lcEmptyOrInvalid || subtractedFilesChanged <= 0
	return changes, emptyOrInvalid
}

// YearlyLineChangeMap
func (ylcm *YearlyLineChangeMap) AddLineChanges(lineChangesToAdd *LineChanges, commitYear int) {
	AdditiveValueMapInsert[int, *LineChanges, YearlyLineChangeMap](*ylcm, commitYear, AddLineChanges, lineChangesToAdd)
}

func (ylcm *YearlyLineChangeMap) SubtractLineChanges(lineChangesToSubtract *LineChanges, commitYear int) {
	SubtractiveValueMapRemove[int, *LineChanges, YearlyLineChangeMap](*ylcm, commitYear, SubtractLineChanges, lineChangesToSubtract)
}

func (ylcm *YearlyLineChangeMap) AddYearlyLineChangeMap(ylcmToAdd YearlyLineChangeMap) {
	for year, lineChangesToAdd := range ylcmToAdd {
		ylcm.AddLineChanges(lineChangesToAdd, year)
	}
}

func (ylcm *YearlyLineChangeMap) SubtractYearlyLineChangeMap(ylcmToSubtract YearlyLineChangeMap) {
	for year, lineChangesToSubtract := range ylcmToSubtract {
		ylcm.SubtractLineChanges(lineChangesToSubtract, year)
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
	AdditiveValueMapInsert[int, *Changes, YearlyChangeMap](*ycm, commitYear, AddChanges, changesToAdd)
}

func (ycm *YearlyChangeMap) SubtractChanges(changesToSubtract *Changes, commitYear int) {
	SubtractiveValueMapRemove[int, *Changes, YearlyChangeMap](*ycm, commitYear, SubtractChanges, changesToSubtract)
}

func (ycm *YearlyChangeMap) LineChanges() YearlyLineChangeMap {
	ylcm := make(YearlyLineChangeMap, 0)

	for year, changes := range *ycm {
		ylcm[year] = &changes.LineChanges
	}

	return ylcm
}
