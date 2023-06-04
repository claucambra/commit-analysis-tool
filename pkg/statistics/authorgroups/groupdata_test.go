package authorgroups

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const expectedInsertCorrel = 0.5337756465
const expectedDeleteCorrel = 0.6832214793
const expectedAuthorCorrel = 0.6715084138

func tenDecimalPlaceRound(num float64) float64 {
	return math.Round(num*10000000000) / 10000000000
}

func TestCorrelation(t *testing.T) {
	groupData := testGroupData()
	unknownGroupData := testUnknownGroupData()

	insertCorrel, deleteCorrel, authorCorrel := groupData.Correlation(unknownGroupData)
	roundInsertCorrel := tenDecimalPlaceRound(insertCorrel)
	roundDeleteCorrel := tenDecimalPlaceRound(deleteCorrel)
	roundAuthorCorrel := tenDecimalPlaceRound(authorCorrel)

	if !cmp.Equal(expectedInsertCorrel, roundInsertCorrel) {
		t.Fatalf("Values did not match: %s", cmp.Diff(expectedInsertCorrel, roundInsertCorrel))
	} else if !cmp.Equal(expectedDeleteCorrel, roundDeleteCorrel) {
		t.Fatalf("Values did not match: %s", cmp.Diff(expectedDeleteCorrel, roundDeleteCorrel))
	} else if !cmp.Equal(expectedAuthorCorrel, roundAuthorCorrel) {
		t.Fatalf("Values did not match: %s", cmp.Diff(expectedAuthorCorrel, roundAuthorCorrel))
	}
}
