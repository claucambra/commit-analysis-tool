package common

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAddLineChanges(t *testing.T) {
	testTotalInsertions := rand.Int()
	testTotalDeletions := rand.Int()
	testChangeAInsertions := rand.Intn(testTotalInsertions)
	testChangeADeletions := rand.Intn(testTotalDeletions)
	testChangeBInsertions := testTotalInsertions - testChangeAInsertions
	testChangeBDeletions := testTotalDeletions - testChangeADeletions

	testChange := &LineChanges{
		NumInsertions: testTotalInsertions,
		NumDeletions:  testTotalDeletions,
	}

	changeA := &LineChanges{
		NumInsertions: testChangeAInsertions,
		NumDeletions:  testChangeADeletions,
	}
	changeB := &LineChanges{
		NumInsertions: testChangeBInsertions,
		NumDeletions:  testChangeBDeletions,
	}

	changeA.AddLineChanges(changeB)
	if !reflect.DeepEqual(changeA, testChange) {
		t.Fatalf(`Added line changes do not match expected changes: 
			Expected %+v
			Received %+v`, testChange, changeA)
	}
}

func TestSubtractLineChanges(t *testing.T) {
	testChangeAInsertions := rand.Int()
	testChangeADeletions := rand.Int()
	testChangeBInsertions := rand.Intn(testChangeAInsertions)
	testChangeBDeletions := rand.Intn(testChangeADeletions)
	testChangeFinalInsertions := testChangeAInsertions - testChangeBInsertions
	testChangeFinalDeletions := testChangeADeletions - testChangeBDeletions

	testChange := &LineChanges{
		NumInsertions: testChangeFinalInsertions,
		NumDeletions:  testChangeFinalDeletions,
	}

	changeA := &LineChanges{
		NumInsertions: testChangeAInsertions,
		NumDeletions:  testChangeADeletions,
	}
	changeB := &LineChanges{
		NumInsertions: testChangeBInsertions,
		NumDeletions:  testChangeBDeletions,
	}

	changeA.SubtractLineChanges(changeB)
	if !reflect.DeepEqual(changeA, testChange) {
		t.Fatalf(`Subtracted line changes do not match expected changes:
			Expected %+v
			Received %+v`, testChange, changeA)
	}
}
