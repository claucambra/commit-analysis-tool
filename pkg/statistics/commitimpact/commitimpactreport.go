package commitimpact

import (
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

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
