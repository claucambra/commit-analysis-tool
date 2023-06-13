package authorgroups

import (
	"log"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics"
)

// This should ideally be an independent struct that we sub-struct with author
// group-specific functionality
type GroupSurvivalReport struct {
	Authors           common.EmailSet
	AuthorsInTimeStep statistics.TimeStepPopulation
	AuthorsSurvival   statistics.TimeStepSurvival

	sqlb *db.SQLiteBackend
}

func NewGroupSurvivalReport(sqlb *db.SQLiteBackend, yearlyAuthors common.EmailSet) *GroupSurvivalReport {
	return &GroupSurvivalReport{
		Authors:           yearlyAuthors,
		AuthorsInTimeStep: statistics.TimeStepPopulation{},
		AuthorsSurvival:   statistics.TimeStepSurvival{},
		sqlb:              sqlb,
	}
}

// See how long each author lasts for on average
func (gsp *GroupSurvivalReport) Generate() {
	gsp.AuthorsInTimeStep = statistics.TimeStepPopulation{}
	gsp.AuthorsSurvival = statistics.TimeStepSurvival{}

	for author := range gsp.Authors {
		timeSteps, err := authorContinuousMonths(gsp.sqlb, author)
		if err != nil {
			log.Printf("Author %s did not have retrievable months", author)
			continue
		} else if timeSteps < 1 {
			log.Printf("Author %s active for less than one month. Can't analyse.", author)
			continue
		}

		for i := 0; i < timeSteps; i++ {
			if len(gsp.AuthorsInTimeStep) < i+1 {
				gsp.AuthorsInTimeStep = append(gsp.AuthorsInTimeStep, 1)
			} else {
				gsp.AuthorsInTimeStep[i] += 1
			}
		}
	}

	gsp.AuthorsSurvival = gsp.AuthorsInTimeStep.KaplanMeierSurvival()
}
