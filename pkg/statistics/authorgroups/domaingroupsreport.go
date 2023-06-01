package authorgroups

import (
	"strings"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

const fallbackDomain = "unknown-domain"
const fallbackGroupName = "unknown"

type DomainGroupsReport struct {
	TotalAuthors           int
	TotalCommits           int
	TotalChanges           *common.LineChanges
	TotalYearlyLineChanges common.YearlyLineChangeMap

	GroupsOfDomains map[string][]string

	DomainTotalNumAuthors   map[string]int
	DomainTotalLineChanges  map[string]*common.LineChanges
	DomainTotalNumDeletions map[string]int

	DomainTotalYearlyLineChanges map[string]common.YearlyLineChangeMap
}

func NewDomainGroupsReport(domainGroups map[string][]string) *DomainGroupsReport {
	report := &DomainGroupsReport{
		TotalChanges:                 &common.LineChanges{},
		TotalYearlyLineChanges:       common.YearlyLineChangeMap{},
		GroupsOfDomains:              domainGroups,
		DomainTotalNumAuthors:        map[string]int{},
		DomainTotalLineChanges:       map[string]*common.LineChanges{},
		DomainTotalYearlyLineChanges: map[string]common.YearlyLineChangeMap{},
	}

	return report
}

func (report *DomainGroupsReport) updateDomainChanges(sqlb *db.SQLiteBackend) {
	for authorDomain := range report.DomainTotalNumAuthors {
		changes, err := domainChanges(sqlb, authorDomain)
		if err != nil {
			return
		}

		report.TotalChanges = common.AddLineChanges(report.TotalChanges, &changes.LineChanges)

		if existingDomainLineChanges, ok := report.DomainTotalLineChanges[authorDomain]; ok {
			summedDomainLineChanges := common.AddLineChanges(existingDomainLineChanges, &changes.LineChanges)
			report.DomainTotalLineChanges[authorDomain] = summedDomainLineChanges
		} else {
			report.DomainTotalLineChanges[authorDomain] = &changes.LineChanges
		}

		yearlyChanges, err := domainYearlyChanges(sqlb, authorDomain)
		if err != nil {
			return
		}

		yearlyLineChanges := yearlyChanges.LineChanges()
		report.TotalYearlyLineChanges.AddYearlyLineChangeMap(yearlyLineChanges)

		if existingDomainYearLineChanges, ok := report.DomainTotalYearlyLineChanges[authorDomain]; ok {
			existingDomainYearLineChanges.AddYearlyLineChangeMap(yearlyLineChanges)
			report.DomainTotalYearlyLineChanges[authorDomain] = existingDomainYearLineChanges
		} else {
			report.DomainTotalYearlyLineChanges[authorDomain] = yearlyLineChanges
		}
	}
}

func (report *DomainGroupsReport) updateAuthors(authors []string, db *db.SQLiteBackend) {
	for _, author := range authors {
		if author == "" {
			continue
		}

		authorDomain := fallbackDomain
		splitAuthorEmail := strings.Split(author, "@")

		if len(splitAuthorEmail) >= 2 {
			authorDomain = splitAuthorEmail[1]
		}

		report.DomainTotalNumAuthors[authorDomain] += 1
		report.TotalAuthors += 1
	}
}

func (report *DomainGroupsReport) Generate(db *db.SQLiteBackend) {
	authors, err := db.Authors()
	if err != nil {
		return
	}

	report.updateAuthors(authors, db)
	report.updateDomainChanges(db)
}

// Returns authors, insertions, deletions
func (report *DomainGroupsReport) accumulateGroupCounts(groupName string) (int, *common.LineChanges, common.YearlyLineChangeMap) {
	totalGroupAuthors := 0
	totalGroupLineChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}
	totalGroupYearlyLineChanges := make(common.YearlyLineChangeMap, 0)

	for _, domain := range report.GroupsOfDomains[groupName] {
		reportChanges, ok := report.DomainTotalLineChanges[domain]
		if !ok {
			continue
		}

		totalGroupLineChanges = common.AddLineChanges(totalGroupLineChanges, reportChanges)
		totalGroupAuthors += report.DomainTotalNumAuthors[domain]

		reportYearlyChanges := report.DomainTotalYearlyLineChanges[domain]
		totalGroupYearlyLineChanges.AddYearlyLineChangeMap(reportYearlyChanges)
	}

	return totalGroupAuthors, totalGroupLineChanges, totalGroupYearlyLineChanges
}

func (report *DomainGroupsReport) UnknownGroupData() *GroupData {
	totalGroupAuthors := 0
	totalGroupChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}
	totalGroupYearlyLineChanges := make(common.YearlyLineChangeMap, 0)

	for groupName := range report.GroupsOfDomains {
		groupAuthors, groupLineChanges, yearlyGroupLineChanges := report.accumulateGroupCounts(groupName)
		totalGroupAuthors += groupAuthors
		totalGroupChanges = common.AddLineChanges(totalGroupChanges, groupLineChanges)
		totalGroupYearlyLineChanges.AddYearlyLineChangeMap(yearlyGroupLineChanges)
	}

	unknownGroupTotalAuthors := report.TotalAuthors - totalGroupAuthors
	unknownGroupTotalLineChanges := common.SubtractLineChanges(report.TotalChanges, totalGroupChanges)
	unknownGroupTotalYearlyLineChanges := report.TotalYearlyLineChanges
	unknownGroupTotalYearlyLineChanges.SubtractYearlyLineChangeMap(totalGroupYearlyLineChanges)

	return NewGroupData(report, fallbackGroupName, unknownGroupTotalAuthors, unknownGroupTotalLineChanges, unknownGroupTotalYearlyLineChanges)
}

func (report *DomainGroupsReport) GroupData(groupName string) *GroupData {
	if groupName == "" || groupName == fallbackGroupName {
		return report.UnknownGroupData()
	}

	totalGroupAuthors, totalGroupLineChanges, totalGroupYearlyLineChanges := report.accumulateGroupCounts(groupName)
	return NewGroupData(report, groupName, totalGroupAuthors, totalGroupLineChanges, totalGroupYearlyLineChanges)
}
