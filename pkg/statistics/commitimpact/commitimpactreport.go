package commitimpact

import (
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

const featureKey = "feature"
const bugfixKey = "bugfix"
const documentationKey = "documentation"
const testingKey = "testing"
const testDataKey = "testdata"

type CommitImpactReport struct {
	Commits    common.CommitMap
	Impact     map[string]float64
	MeanImpact float64
}

func NewCommitImpactReport(commits common.CommitMap) *CommitImpactReport {
	return &CommitImpactReport{
		Commits: commits,
		Impact:  map[string]float64{},
	}
}

func codingWeightMap() map[string]float64 {
	return map[string]float64{
		featureKey:       1.0,
		bugfixKey:        0.8,
		documentationKey: 0.6,
		testingKey:       0.3,
		testDataKey:      0.0,
	}
}
