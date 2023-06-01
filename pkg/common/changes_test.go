package common

import (
	"math/rand"
	"reflect"
	"testing"
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

	changeA.AddLineChanges(changeB)
	if !reflect.DeepEqual(changeA, testChange) {
		t.Fatalf(`Added line changes do not match expected changes: 
			Expected %+v
			Received %+v`, testChange, changeA)
	}
}

func TestSubtractLineChanges(t *testing.T) {
	changeA, changeB := generateRandomLineChanges()
	testChange := &LineChanges{
		NumInsertions: changeA.NumInsertions - changeB.NumInsertions,
		NumDeletions:  changeA.NumDeletions - changeB.NumDeletions,
	}

	changeA.SubtractLineChanges(changeB)
	if !reflect.DeepEqual(changeA, testChange) {
		t.Fatalf(`Subtracted line changes do not match expected changes:
			Expected %+v
			Received %+v`, testChange, changeA)
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

	changeA.AddChanges(changeB)
	if !reflect.DeepEqual(changeA, testChange) {
		t.Fatalf(`Added changes do not match expected changes:
			Expected %+v
			Received %+v`, testChange, changeA)
	}
}

func TestSubtractChanges(t *testing.T) {
	changeA, changeB := generateRandomChanges()
	testChange := &Changes{
		LineChanges: LineChanges{
			NumInsertions: changeA.NumInsertions - changeB.NumInsertions,
			NumDeletions:  changeA.NumDeletions - changeB.NumDeletions,
		},
		NumFilesChanged: changeA.NumFilesChanged - changeB.NumFilesChanged,
	}

	changeA.SubtractChanges(changeB)
	if !reflect.DeepEqual(changeA, testChange) {
		t.Fatalf(`Subtracted changes do not match expected changes:
			Expected %+v
			Received %+v`, testChange, changeA)
	}
}

// YearlyLineChangeMap
func TestAddLineChangesInYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYear := 2023
	ylcm := make(YearlyLineChangeMap, 0)

	ylcm.AddLineChanges(lineChangeA, testYear)
	if ylcmALineChangeA := ylcm[testYear]; !reflect.DeepEqual(ylcmALineChangeA, *lineChangeA) {
		t.Fatalf(`Added changes to yearly line change map when year not already in map does not match expected changes:
			Expected %+v
			Received %+v`, lineChangeA, ylcmALineChangeA)
	}

	ylcm.AddLineChanges(lineChangeB, testYear)
	summedLineChanges := *lineChangeA
	summedLineChanges.AddLineChanges(lineChangeB)

	if ylcmASummedLineChange := ylcm[testYear]; !reflect.DeepEqual(ylcmASummedLineChange, summedLineChanges) {
		t.Fatalf(`Added changes to yearly line change map when year already in map does not match expected changes:
			Expected %+v
			Received %+v`, summedLineChanges, ylcmASummedLineChange)
	}
}

func TestSubtractLineChangesInYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2023
	testYearB := 2003
	ylcm := YearlyLineChangeMap{testYearA: *lineChangeA}

	ylcm.SubtractLineChanges(lineChangeB, testYearB)
	if _, ok := ylcm[testYearB]; ok {
		t.Fatalf("Subtracting line changes from a year not present in YLCM should not add this year to YLCM.")
	}

	ylcm = YearlyLineChangeMap{testYearA: *lineChangeA}
	ylcm.SubtractLineChanges(lineChangeB, testYearA)

	expectedSubLineChanges := LineChanges{
		NumInsertions: lineChangeA.NumInsertions - lineChangeB.NumInsertions,
		NumDeletions:  lineChangeA.NumDeletions - lineChangeB.NumDeletions,
	}

	if subLineChanges := ylcm[testYearA]; !reflect.DeepEqual(subLineChanges, expectedSubLineChanges) {
		t.Fatalf(`Subtracted line changes from yearly line change map does not match expected changes:
			Expected %+v
			Received %+v`, expectedSubLineChanges, subLineChanges)
	}
}

func TestAddYearlyLineChangeMapToYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2023
	testYearB := 2004
	ylcmA := YearlyLineChangeMap{testYearA: *lineChangeA}
	ylcmB := YearlyLineChangeMap{testYearB: *lineChangeB}

	ylcmA.AddYearlyLineChangeMap(ylcmB)

	if addedLineChangeB, ok := ylcmA[testYearB]; !ok {
		t.Fatalf("Adding line changes from a year of YLCM B not present in YLCM A should add this year to YLCM A.")
	} else if !reflect.DeepEqual(addedLineChangeB, *lineChangeB) {
		t.Fatalf(`Added YLCM B to YLCM A does not match expected line changes:
			Expected %+v
			Received %+v`, *lineChangeB, addedLineChangeB)
	}

	ylcmB = YearlyLineChangeMap{testYearA: *lineChangeB}
	ylcmA.AddYearlyLineChangeMap(ylcmB)

	expectedAddLineChanges := LineChanges{
		NumInsertions: lineChangeA.NumInsertions + lineChangeB.NumInsertions,
		NumDeletions:  lineChangeA.NumDeletions + lineChangeB.NumDeletions,
	}

	if addedLineChange := ylcmA[testYearA]; !reflect.DeepEqual(addedLineChange, expectedAddLineChanges) {
		t.Fatalf(`Added YLCM B line changes to YLCM A does not match expected line changes:
			Expected %+v
			Received %+v`, addedLineChange, expectedAddLineChanges)
	}
}

func TestSubtractYearlyLineChangeMapToYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2023
	testYearB := 2004
	ylcmA := YearlyLineChangeMap{testYearA: *lineChangeA}
	ylcmB := YearlyLineChangeMap{testYearB: *lineChangeB}

	ylcmA.SubtractYearlyLineChangeMap(ylcmB)

	if _, ok := ylcmA[testYearB]; ok {
		t.Fatalf("Subtracting YLCM B line changes from a year not present in YLCM A should not add this year to YLCM A.")
	}

	ylcmA.AddLineChanges(lineChangeA, testYearB)
	ylcmA.SubtractYearlyLineChangeMap(ylcmB)

	expectedSubLineChanges := LineChanges{
		NumInsertions: lineChangeA.NumInsertions - lineChangeB.NumInsertions,
		NumDeletions:  lineChangeA.NumDeletions - lineChangeB.NumDeletions,
	}

	if testYearBSubLineChanges, ok := ylcmA[testYearB]; !ok {
		t.Fatalf("Test year B should now be present in YLCM A")
	} else if !reflect.DeepEqual(testYearBSubLineChanges, expectedSubLineChanges) {
		t.Fatalf(`Subtracted YLCM B line changes to yearly line change map does not match expected changes:
			Expected %+v
			Received %+v`, expectedSubLineChanges, testYearBSubLineChanges)
	}
}

func TestSeparatedChangeArrayFromYearlyLineChangeMap(t *testing.T) {
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2001
	testYearB := 2023
	testYearC := 2195
	ylcm := YearlyLineChangeMap{
		testYearA: *lineChangeA,
		testYearB: *lineChangeB,
		testYearC: *lineChangeA,
	}

	expectedInsertionsArray := []int{lineChangeA.NumInsertions, lineChangeB.NumInsertions, lineChangeA.NumInsertions}
	expectedDeletionsArray := []int{lineChangeA.NumDeletions, lineChangeB.NumDeletions, lineChangeA.NumDeletions}

	insertionsArray, deletionsArray := ylcm.SeparatedChangeArrays(nil)

	if !reflect.DeepEqual(insertionsArray, expectedInsertionsArray) {
		t.Fatalf(`Insertions array from yearly change map does not match expected insertions array:
			Expected %+v
			Received %+v`, expectedInsertionsArray, insertionsArray)
	}

	if !reflect.DeepEqual(deletionsArray, expectedDeletionsArray) {
		t.Fatalf(`Deletions array from yearly change map does not match expected deletions array:
			Expected %+v
			Received %+v`, expectedDeletionsArray, deletionsArray)
	}

	testRetrievalYears := []int{testYearA, testYearC}
	specificInsertionsArray, specificDeletionsArray := ylcm.SeparatedChangeArrays(testRetrievalYears)
	expectedSpecificInsertionsArray := []int{lineChangeA.NumInsertions, lineChangeA.NumInsertions}
	expectedSpecificDeletionsArray := []int{lineChangeA.NumDeletions, lineChangeA.NumDeletions}

	if !reflect.DeepEqual(specificInsertionsArray, expectedSpecificInsertionsArray) {
		t.Fatalf(`Insertions array from specific years in yearly change map does not match expected insertions array:
			Expected %+v
			Received %+v`, expectedSpecificInsertionsArray, specificInsertionsArray)
	}

	if !reflect.DeepEqual(specificDeletionsArray, expectedSpecificDeletionsArray) {
		t.Fatalf(`Deletions array from specific yearly change map does not match expected deletions array:
			Expected %+v
			Received %+v`, expectedSpecificDeletionsArray, specificDeletionsArray)
	}
}

// YearlyChangeMap
func TestAddChangesToYearlyChangeMap(t *testing.T) {
	changeA, changeB := generateRandomChanges()
	testYear := 2023
	ycm := make(YearlyChangeMap, 0)

	ycm.AddChanges(changeA, testYear)
	if ycmChangeA := ycm[testYear]; !reflect.DeepEqual(ycmChangeA, *changeA) {
		t.Fatalf(`Added changes to yearly change map when year not already in map does not match expected changes:
			Expected %+v
			Received %+v`, *changeA, ycmChangeA)
	}

	ycm.AddChanges(changeB, testYear)
	summedChanges := *changeA
	summedChanges.AddChanges(changeB)

	if ycmASummedChange := ycm[testYear]; !reflect.DeepEqual(ycmASummedChange, summedChanges) {
		t.Fatalf(`Added changes to yearly change map when year already in map does not match expected changes:
			Expected %+v
			Received %+v`, summedChanges, ycmASummedChange)
	}
}

func TestSubtractChangesFromYearlyChangeMap(t *testing.T) {
	changeA, changeB := generateRandomChanges()
	testYearA := 2023
	testYearB := 2003
	ycm := YearlyChangeMap{testYearA: *changeA}

	ycm.SubtractChanges(changeB, testYearB)
	if _, ok := ycm[testYearB]; ok {
		t.Fatalf("Subtracting changes from a year not present in YCM should not add this year to YCM.")
	}

	ycm.SubtractChanges(changeB, testYearA)

	expectedSubChanges := Changes{
		LineChanges: LineChanges{
			NumInsertions: changeA.NumInsertions - changeB.NumInsertions,
			NumDeletions:  changeA.NumDeletions - changeB.NumDeletions,
		},
		NumFilesChanged: changeA.NumFilesChanged - changeB.NumFilesChanged,
	}

	if subChanges := ycm[testYearA]; !reflect.DeepEqual(subChanges, expectedSubChanges) {
		t.Fatalf(`Subtracted line changes from yearly line change map does not match expected changes:
			Expected %+v
			Received %+v`, expectedSubChanges, subChanges)
	}
}
