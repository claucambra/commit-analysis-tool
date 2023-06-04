package common

import (
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const randomChangeMaxLimit = 50000

func generateRandomLineChanges() (*LineChanges, *LineChanges) {
	testChangeAInsertions := rand.Intn(randomChangeMaxLimit)
	testChangeADeletions := rand.Intn(randomChangeMaxLimit)
	testChangeBInsertions := rand.Intn(testChangeAInsertions)
	testChangeBDeletions := rand.Intn(testChangeADeletions)

	changeA := &LineChanges{
		NumInsertions: testChangeAInsertions,
		NumDeletions:  testChangeADeletions,
	}
	changeB := &LineChanges{
		NumInsertions: testChangeBInsertions,
		NumDeletions:  testChangeBDeletions,
	}

	return changeA, changeB
}

func generateRandomChanges() (*Changes, *Changes) {
	changeAFilesChanged := rand.Intn(randomChangeMaxLimit)
	changeBFilesChanged := rand.Intn(randomChangeMaxLimit)
	changesALineChanges, changeBLineChanges := generateRandomLineChanges()

	changeA := &Changes{
		LineChanges:     *changesALineChanges,
		NumFilesChanged: changeAFilesChanged,
	}

	changeB := &Changes{
		LineChanges:     *changeBLineChanges,
		NumFilesChanged: changeBFilesChanged,
	}

	return changeA, changeB
}

// LineChanges
func TestAddLineChanges(t *testing.T) {
	changeA, changeB := generateRandomLineChanges()
	testChange := &LineChanges{
		NumInsertions: changeA.NumInsertions + changeB.NumInsertions,
		NumDeletions:  changeA.NumDeletions + changeB.NumDeletions,
	}

	summedChanges := AddLineChanges(changeA, changeB)
	if !cmp.Equal(summedChanges, testChange) {
		t.Fatalf(`Added line changes do not match expected changes: %s`, cmp.Diff(testChange, summedChanges))
	}
}

func TestSubtractLineChanges(t *testing.T) {
	changeA, changeB := generateRandomLineChanges()
	testChange := &LineChanges{
		NumInsertions: MaxInt(changeA.NumInsertions-changeB.NumInsertions, 0),
		NumDeletions:  MaxInt(changeA.NumDeletions-changeB.NumDeletions, 0),
	}

	subbedChanges, _ := SubtractLineChanges(changeA, changeB)
	if !cmp.Equal(subbedChanges, testChange) {
		t.Fatalf(`Subtracted line changes do not match expected changes: %s`, cmp.Diff(testChange, subbedChanges))
	}
}

// Changes
func TestAddChanges(t *testing.T) {
	changeA, changeB := generateRandomChanges()
	testChange := &Changes{
		LineChanges: LineChanges{
			NumInsertions: changeA.NumInsertions + changeB.NumInsertions,
			NumDeletions:  changeA.NumDeletions + changeB.NumDeletions,
		},
		NumFilesChanged: changeA.NumFilesChanged + changeB.NumFilesChanged,
	}

	summedChanges := AddChanges(changeA, changeB)
	if !cmp.Equal(summedChanges, testChange) {
		t.Fatalf(`Added changes do not match expected changes: %s`, cmp.Diff(testChange, summedChanges))
	}
}

func TestSubtractChanges(t *testing.T) {
	changeA, changeB := generateRandomChanges()
	testChange := &Changes{
		LineChanges: LineChanges{
			NumInsertions: MaxInt(changeA.NumInsertions-changeB.NumInsertions, 0),
			NumDeletions:  MaxInt(changeA.NumDeletions-changeB.NumDeletions, 0),
		},
		NumFilesChanged: MaxInt(changeA.NumFilesChanged-changeB.NumFilesChanged, 0),
	}

	subbedChanges, _ := SubtractChanges(changeA, changeB)
	if !cmp.Equal(subbedChanges, testChange) {
		t.Fatalf(`Subtracted changes do not match expected changes: %s`, cmp.Diff(testChange, subbedChanges))
	}
}

// YearlyLineChangeMap
func TestAddLineChangesInYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYear := 2023
	ylcm := make(YearlyLineChangeMap, 0)

	ylcm.AddLineChanges(lineChangeA, testYear)
	if ylcmALineChangeA := ylcm[testYear]; !cmp.Equal(ylcmALineChangeA, lineChangeA) {
		t.Fatalf(`Added changes to yearly line change map when year not already in map does not match expected changes: %s`, cmp.Diff(lineChangeA, ylcmALineChangeA))
	}

	ylcm.AddLineChanges(lineChangeB, testYear)
	summedLineChanges := AddLineChanges(lineChangeA, lineChangeB)

	if ylcmASummedLineChange := ylcm[testYear]; !cmp.Equal(ylcmASummedLineChange, summedLineChanges) {
		t.Fatalf(`Added changes to yearly line change map when year already in map does not match expected changes: %s`, cmp.Diff(summedLineChanges, ylcmASummedLineChange))
	}
}

func TestSubtractLineChangesInYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2023
	testYearB := 2003
	ylcm := YearlyLineChangeMap{testYearA: lineChangeA}

	ylcm.SubtractLineChanges(lineChangeB, testYearB)
	if _, ok := ylcm[testYearB]; ok {
		t.Fatalf("Subtracting line changes from a year not present in YLCM should not add this year to YLCM.")
	}

	ylcm = YearlyLineChangeMap{testYearA: lineChangeA}
	ylcm.SubtractLineChanges(lineChangeB, testYearA)

	expectedSubLineChanges := &LineChanges{
		NumInsertions: MaxInt(lineChangeA.NumInsertions-lineChangeB.NumInsertions, 0),
		NumDeletions:  MaxInt(lineChangeA.NumDeletions-lineChangeB.NumDeletions, 0),
	}

	if subLineChanges := ylcm[testYearA]; !cmp.Equal(subLineChanges, expectedSubLineChanges) {
		t.Fatalf(`Subtracted line changes from yearly line change map does not match expected changes: %s`, cmp.Diff(expectedSubLineChanges, subLineChanges))
	}
}

func TestAddYearlyLineChangeMapToYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2023
	testYearB := 2004
	ylcmA := YearlyLineChangeMap{testYearA: lineChangeA}
	ylcmB := YearlyLineChangeMap{testYearB: lineChangeB}

	ylcmA.AddYearlyLineChangeMap(ylcmB)

	if addedLineChangeB, ok := ylcmA[testYearB]; !ok {
		t.Fatalf("Adding line changes from a year of YLCM B not present in YLCM A should add this year to YLCM A.")
	} else if !cmp.Equal(addedLineChangeB, lineChangeB) {
		t.Fatalf(`Added YLCM B to YLCM A does not match expected line changes: %s`, cmp.Diff(*lineChangeB, addedLineChangeB))
	}

	ylcmB = YearlyLineChangeMap{testYearA: lineChangeB}
	ylcmA.AddYearlyLineChangeMap(ylcmB)

	expectedAddLineChanges := &LineChanges{
		NumInsertions: lineChangeA.NumInsertions + lineChangeB.NumInsertions,
		NumDeletions:  lineChangeA.NumDeletions + lineChangeB.NumDeletions,
	}

	if addedLineChange := ylcmA[testYearA]; !cmp.Equal(addedLineChange, expectedAddLineChanges) {
		t.Fatalf(`Added YLCM B line changes to YLCM A does not match expected line changes: %s`, cmp.Diff(addedLineChange, expectedAddLineChanges))
	}
}

func TestSubtractYearlyLineChangeMapToYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2023
	testYearB := 2004
	ylcmA := YearlyLineChangeMap{testYearA: lineChangeA}
	ylcmB := YearlyLineChangeMap{testYearB: lineChangeB}

	ylcmA.SubtractYearlyLineChangeMap(ylcmB)

	if _, ok := ylcmA[testYearB]; ok {
		t.Fatalf("Subtracting YLCM B line changes from a year not present in YLCM A should not add this year to YLCM A.")
	}

	ylcmA.AddLineChanges(lineChangeA, testYearB)
	ylcmA.SubtractYearlyLineChangeMap(ylcmB)

	expectedSubLineChanges := &LineChanges{
		NumInsertions: MaxInt(lineChangeA.NumInsertions-lineChangeB.NumInsertions, 0),
		NumDeletions:  MaxInt(lineChangeA.NumDeletions-lineChangeB.NumDeletions, 0),
	}

	if testYearBSubLineChanges, ok := ylcmA[testYearB]; !ok {
		t.Fatalf("Test year B should now be present in YLCM A")
	} else if !cmp.Equal(testYearBSubLineChanges, expectedSubLineChanges) {
		t.Fatalf(`Subtracted YLCM B line changes to yearly line change map does not match expected changes: %s`, cmp.Diff(expectedSubLineChanges, testYearBSubLineChanges))
	}
}

func TestSeparatedChangeArrayFromYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2001
	testYearB := 2023
	testYearC := 2195
	ylcm := YearlyLineChangeMap{
		testYearA: lineChangeA,
		testYearB: lineChangeB,
		testYearC: lineChangeA,
	}

	expectedInsertionsArray := []int{lineChangeA.NumInsertions, lineChangeB.NumInsertions, lineChangeA.NumInsertions}
	expectedDeletionsArray := []int{lineChangeA.NumDeletions, lineChangeB.NumDeletions, lineChangeA.NumDeletions}

	insertionsArray, deletionsArray := ylcm.SeparatedChangeArrays(nil)

	if !cmp.Equal(insertionsArray, expectedInsertionsArray) {
		t.Fatalf(`Insertions array from yearly change map does not match expected insertions array: %s`, cmp.Diff(expectedInsertionsArray, insertionsArray))
	}

	if !cmp.Equal(deletionsArray, expectedDeletionsArray) {
		t.Fatalf(`Deletions array from yearly change map does not match expected deletions array: %s`, cmp.Diff(expectedDeletionsArray, deletionsArray))
	}

	testRetrievalYears := []int{testYearA, testYearC}
	specificInsertionsArray, specificDeletionsArray := ylcm.SeparatedChangeArrays(testRetrievalYears)
	expectedSpecificInsertionsArray := []int{lineChangeA.NumInsertions, lineChangeA.NumInsertions}
	expectedSpecificDeletionsArray := []int{lineChangeA.NumDeletions, lineChangeA.NumDeletions}

	if !cmp.Equal(specificInsertionsArray, expectedSpecificInsertionsArray) {
		t.Fatalf(`Insertions array from specific years in yearly change map does not match expected insertions array: %s`, cmp.Diff(expectedSpecificInsertionsArray, specificInsertionsArray))
	}

	if !cmp.Equal(specificDeletionsArray, expectedSpecificDeletionsArray) {
		t.Fatalf(`Deletions array from specific yearly change map does not match expected deletions array: %s`, cmp.Diff(expectedSpecificDeletionsArray, specificDeletionsArray))
	}
}

// YearlyChangeMap
func TestAddChangesToYearlyChangeMap(t *testing.T) {
	changeA, changeB := generateRandomChanges()
	testYear := 2023
	ycm := make(YearlyChangeMap, 0)

	ycm.AddChanges(changeA, testYear)
	if ycmChangeA := ycm[testYear]; !cmp.Equal(ycmChangeA, changeA) {
		t.Fatalf(`Added changes to yearly change map when year not already in map does not match expected changes: %s`, cmp.Diff(changeA, ycmChangeA))
	}

	ycm.AddChanges(changeB, testYear)
	summedChanges := AddChanges(changeA, changeB)

	if ycmASummedChange := ycm[testYear]; !cmp.Equal(ycmASummedChange, summedChanges) {
		t.Fatalf(`Added changes to yearly change map when year already in map does not match expected changes: %s`, cmp.Diff(summedChanges, ycmASummedChange))
	}
}

func TestSubtractChangesFromYearlyChangeMap(t *testing.T) {
	changeA, changeB := generateRandomChanges()
	testYearA := 2023
	testYearB := 2003
	ycm := YearlyChangeMap{testYearA: changeA}

	ycm.SubtractChanges(changeB, testYearB)
	if _, ok := ycm[testYearB]; ok {
		t.Fatalf("Subtracting changes from a year not present in YCM should not add this year to YCM.")
	}

	ycm.SubtractChanges(changeB, testYearA)

	expectedSubChanges := &Changes{
		LineChanges: LineChanges{
			NumInsertions: MaxInt(changeA.NumInsertions-changeB.NumInsertions, 0),
			NumDeletions:  MaxInt(changeA.NumDeletions-changeB.NumDeletions, 0),
		},
		NumFilesChanged: MaxInt(changeA.NumFilesChanged-changeB.NumFilesChanged, 0),
	}

	if subChanges := ycm[testYearA]; !cmp.Equal(subChanges, expectedSubChanges) {
		t.Fatalf(`Subtracted line changes from yearly line change map does not match expected changes: %s`, cmp.Diff(expectedSubChanges, subChanges))
	}
}
