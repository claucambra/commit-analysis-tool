package commitimpact

import (
	"log"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/commitcoding"
	"gonum.org/v1/gonum/stat"
)

const featureKey = "feature"
const bugfixKey = "bugfix"
const documentationKey = "documentation"
const testingKey = "testing"
const testDataKey = "testdata"

const insertionWeight = 0.9
const deletionWeight = 0.7

const suspiciouslyHighImpactThreshold = 5000

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

func codeMap() map[string][]string {
	return map[string][]string{
		featureKey:       {`\bintroduc(e|tion)+\b`, `add(ed|ition)*\b\s+(\ba\b)*\s*\b(support|new|option|way|function)(s)*\b`},
		bugfixKey:        {`\bfix(ed|es)*\b`, `\bsanitise\b`, `\bbroken\b`, `\bbreak(s|ing)+\b`, `\brevert(s|ing)*\b`, `add(ed|ition)*\b\s+(\ba\b)*\s*\b(missing)*\b`},
		documentationKey: {`\bdocument\b`, `\bexplain\b`, `\bcomment\b`},
		testingKey:       {`\btest(ing)*\b`},
		testDataKey:      {`\btest data\b`},
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

func (cir *CommitImpactReport) generateImpacts(codeMatchCommits map[string][]*common.Commit) {
	log.Printf("Generating commit impact scores.")

	codeWeightMap := codingWeightMap()
	commitWeights := map[string]float64{}
	commitImpacts := []float64{}

	for codeCategory, commits := range codeMatchCommits {
		for _, commit := range commits {
			commitWeights[commit.Id] = codeWeightMap[codeCategory]
		}
	}

	for commitId, weight := range commitWeights {
		commit, ok := cir.Commits[commitId]
		if !ok {
			log.Fatalf("Could not find commit with id %s in commit impact report commits.", commitId)
			continue
		}

		insertScore := float64(commit.NumInsertions) * insertionWeight
		deleteScore := float64(commit.NumDeletions) * deletionWeight

		impactScore := (insertScore + deleteScore) * weight

		if impactScore > suspiciouslyHighImpactThreshold {
			log.Printf("Found a commit (%s) with a suspiciously high impact score, ignoring.", commitId)
			continue
		}

		// Bias towards lower impact
		commitImpacts = append(commitImpacts, impactScore)
		cir.Impact[commitId] = impactScore
	}

	cir.MeanImpact = stat.Mean(commitImpacts, nil)

	log.Printf("Analysed %v commits, produced a mean impact score of %f", len(commitImpacts), cir.MeanImpact)
}

// Not all commits we have will get impact scores, this depends on the CommitCodingReport
func (cir *CommitImpactReport) Generate() {
	codingReport := commitcoding.NewCommitCodingReport(cir.Commits, codeMap())
	codingReport.Generate()

	cir.generateImpacts(codingReport.CodeMatchCommits)
}
