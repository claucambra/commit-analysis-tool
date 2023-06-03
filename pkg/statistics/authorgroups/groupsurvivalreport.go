package authorgroups

import (
	"log"
	"math"

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

// See how long each author lasts for on average
func (gsp *GroupSurvivalReport) Generate() {
	gsp.AuthorsInYear = statistics.TimeStepPopulation{}
	gsp.AuthorsSurvival = statistics.TimeStepSurvival{}

	for author := range gsp.Authors {
		years, err := authorYears(gsp.sqlb, author)
		if err != nil {
			log.Printf("Author %s did not have retrievable years", author)
			continue
		} else if len(years) < 2 {
			log.Printf("Author %s active for less than two years. Can't analyse.", author)
			continue
		}

		// We only track as long as the developer survives the first time they contribute
		prevYear := math.MinInt
		for i, year := range years {
			if prevYear != math.MinInt && prevYear != year-1 {
				break
			}

			if len(gsp.AuthorsInYear) < i+1 {
				gsp.AuthorsInYear = append(gsp.AuthorsInYear, 1)
			} else {
				gsp.AuthorsInYear[i] += 1
			}

			prevYear = year
		}
	}

	gsp.AuthorsSurvival = gsp.AuthorsInYear.KaplanMeierSurvival()
}
