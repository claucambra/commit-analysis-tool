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
