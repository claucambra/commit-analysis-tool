package common

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

const randomChangeMaxLimit = (math.MaxInt / 2) - 1

func generateRandomLineChanges() (*LineChanges, *LineChanges) {
	testChangeAInsertions := rand.Intn(randomChangeMaxLimit)
	testChangeADeletions := rand.Intn(randomChangeMaxLimit)
	testChangeBInsertions := rand.Intn(randomChangeMaxLimit)
	testChangeBDeletions := rand.Intn(randomChangeMaxLimit)

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

func TestAddYearlyLineChangeMapToYearlyLineChangeMap(t *testing.T) {
	ylcmA := make(YearlyLineChangeMap, 0)
	ylcmB := make(YearlyLineChangeMap, 0)
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2023
	testYearB := 2004

	ylcmA.AddLineChanges(lineChangeA, testYearA)
	ylcmB.AddLineChanges(lineChangeB, testYearB)

	ylcmA.AddYearlyLineChangeMap(ylcmB)

	if addedLineChangeB, ok := ylcmA[testYearB]; !ok {
		t.Fatalf("Adding line changes from a year of YLCM B not present in YLCM A should add this year to YLCM A.")
	} else if !reflect.DeepEqual(addedLineChangeB, *lineChangeB) {
		t.Fatalf(`Added YLCM B to YLCM A does not match expected line changes:
			Expected %+v
			Received %+v`, *lineChangeB, addedLineChangeB)
	}

	ylcmB.AddLineChanges(lineChangeB, testYearA)
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
	ylcmA := make(YearlyLineChangeMap, 0)
	ylcmB := make(YearlyLineChangeMap, 0)
	lineChangeA, lineChangeB := generateRandomLineChanges()
	testYearA := 2023
	testYearB := 2004

	ylcmA.AddLineChanges(lineChangeA, testYearA)
	ylcmB.AddLineChanges(lineChangeB, testYearB)

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

func TestAddChangesToYearlyChangeMap(t *testing.T) {
	changeA, changeB := generateRandomChanges()
	testYear := 2023
	ycm := make(YearlyChangeMap, 0)

	ycm.AddChanges(changeA, testYear)
	if ycmAChangeA := *(ycm[testYear]); !reflect.DeepEqual(ycmAChangeA, *changeA) {
		t.Fatalf(`Added changes to yearly change map when year not already in map does not match expected changes:
			Expected %+v
			Received %+v`, *changeA, ycmAChangeA)
	}

	ycm.AddChanges(changeB, testYear)
	summedChanges := *changeA
	summedChanges.AddChanges(changeB)

	if ycmASummedChange := *(ycm[testYear]); !reflect.DeepEqual(ycmASummedChange, summedChanges) {
		t.Fatalf(`Added changes to yearly change map when year already in map does not match expected changes:
			Expected %+v
			Received %+v`, summedChanges, ycmASummedChange)
	}
}
