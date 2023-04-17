package statistics

import "github.com/claucambra/commit-analysis-tool/internal/git"

type StatisticsReport interface {
	AddCommit(git.CommitData)
	String()
}
