package statistics

import (
	"github.com/claucambra/commit-analysis-tool/internal/git"
)

type CorpAuthorsReport struct {
	Commits            []*git.CommitData
	TotalAuthors       int
	NumCorpAuthors     int
	CorpAuthorsPercent float32
	DomainCountMap     map[string]int
}

