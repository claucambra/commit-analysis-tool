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
