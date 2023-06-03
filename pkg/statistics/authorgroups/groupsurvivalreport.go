package authorgroups

import (
	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics"
)

// This should ideally be an independent struct that we sub-struct with author
// group-specific functionality
type GroupSurvivalReport struct {
	Authors         common.EmailSet
	AuthorsInYear   statistics.TimeStepPopulation
	AuthorsSurvival statistics.TimeStepSurvival
	sqlb            *db.SQLiteBackend
}

func NewGroupSurvivalReport(sqlb *db.SQLiteBackend, yearlyAuthors common.EmailSet) *GroupSurvivalReport {
	return &GroupSurvivalReport{
		Authors:         yearlyAuthors,
		AuthorsInYear:   statistics.TimeStepPopulation{},
		AuthorsSurvival: statistics.TimeStepSurvival{},
		sqlb:            sqlb,
	}
}
