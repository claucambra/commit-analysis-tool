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
	TotalYearlyAuthors     common.YearlyEmailMap
	TotalYearlyLineChanges common.YearlyLineChangeMap

	GroupsOfDomains map[string][]string

	DomainTotalNumAuthors   map[string]int
	DomainTotalLineChanges  map[string]*common.LineChanges
	DomainTotalNumDeletions map[string]int

	DomainTotalYearlyAuthors     map[string]common.YearlyEmailMap
	DomainTotalYearlyLineChanges map[string]common.YearlyLineChangeMap
}

func NewDomainGroupsReport(domainGroups map[string][]string) *DomainGroupsReport {
	return &DomainGroupsReport{
		TotalChanges:                 &common.LineChanges{},
		TotalYearlyAuthors:           common.YearlyEmailMap{},
		TotalYearlyLineChanges:       common.YearlyLineChangeMap{},
		GroupsOfDomains:              domainGroups,
		DomainTotalNumAuthors:        map[string]int{},
		DomainTotalLineChanges:       map[string]*common.LineChanges{},
		DomainTotalYearlyAuthors:     map[string]common.YearlyEmailMap{},
		DomainTotalYearlyLineChanges: map[string]common.YearlyLineChangeMap{},
	}
}

func (report *DomainGroupsReport) resetStats() {
	report.TotalChanges = &common.LineChanges{}
	report.TotalYearlyLineChanges = common.YearlyLineChangeMap{}
	report.DomainTotalNumAuthors = map[string]int{}
	report.DomainTotalLineChanges = map[string]*common.LineChanges{}
	report.DomainTotalYearlyLineChanges = map[string]common.YearlyLineChangeMap{}
}

func (report *DomainGroupsReport) updateDomainChanges(sqlb *db.SQLiteBackend) {
	for authorDomain := range report.DomainTotalNumAuthors {
		lineChanges, err := domainLineChanges(sqlb, authorDomain)
		if err != nil {
			return
		}

		report.TotalChanges = common.AddLineChanges(report.TotalChanges, lineChanges)

		if existingDomainLineChanges, ok := report.DomainTotalLineChanges[authorDomain]; ok {
			summedDomainLineChanges := common.AddLineChanges(existingDomainLineChanges, lineChanges)
			report.DomainTotalLineChanges[authorDomain] = summedDomainLineChanges
		} else {
			report.DomainTotalLineChanges[authorDomain] = lineChanges
		}

		yearlyLineChanges, err := domainYearlyLineChanges(sqlb, authorDomain)
		if err != nil {
			return
		}

		report.TotalYearlyLineChanges.AddYearlyLineChangeMap(yearlyLineChanges)

		if existingDomainYearLineChanges, ok := report.DomainTotalYearlyLineChanges[authorDomain]; ok {
			existingDomainYearLineChanges.AddYearlyLineChangeMap(yearlyLineChanges)
			report.DomainTotalYearlyLineChanges[authorDomain] = existingDomainYearLineChanges
		} else {
			report.DomainTotalYearlyLineChanges[authorDomain] = yearlyLineChanges
		}

		yearlyAuthors, err := domainYearlyAuthors(sqlb, authorDomain)
		if err != nil {
			return
		}

		report.TotalYearlyAuthors.AddYearlyEmailMap(yearlyAuthors)

		if existingDomainYearAuthors, ok := report.DomainTotalYearlyAuthors[authorDomain]; ok {
			existingDomainYearAuthors.AddYearlyEmailMap(yearlyAuthors)
			report.DomainTotalYearlyAuthors[authorDomain] = existingDomainYearAuthors
		} else {
			report.DomainTotalYearlyAuthors[authorDomain] = yearlyAuthors
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

	report.resetStats()
	report.updateAuthors(authors, db)
	report.updateDomainChanges(db)
}

// Returns authors, insertions, deletions
func (report *DomainGroupsReport) accumulateGroupCounts(groupName string) (int, *common.LineChanges, common.YearlyEmailMap, common.YearlyLineChangeMap) {
	totalGroupAuthors := 0
	totalGroupLineChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}
	totalGroupYearlyLineChanges := make(common.YearlyLineChangeMap, 0)
	totalGroupYearlyAuthors := make(common.YearlyEmailMap, 0)

	for _, domain := range report.GroupsOfDomains[groupName] {
		reportChanges, ok := report.DomainTotalLineChanges[domain]
		if !ok {
			continue
		}

		totalGroupLineChanges = common.AddLineChanges(totalGroupLineChanges, reportChanges)
		totalGroupAuthors += report.DomainTotalNumAuthors[domain]

		reportYearlyChanges := report.DomainTotalYearlyLineChanges[domain]
		totalGroupYearlyLineChanges.AddYearlyLineChangeMap(reportYearlyChanges)

		reportYearlyAuthors := report.DomainTotalYearlyAuthors[domain]
		totalGroupYearlyAuthors.AddYearlyEmailMap(reportYearlyAuthors)
	}

	return totalGroupAuthors, totalGroupLineChanges, totalGroupYearlyAuthors, totalGroupYearlyLineChanges
}

func (report *DomainGroupsReport) UnknownGroupData() *GroupData {
	totalGroupAuthors := 0
	totalGroupChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}
	totalGroupYearlyLineChanges := make(common.YearlyLineChangeMap, 0)
	totalGroupYearlyAuthors := make(common.YearlyEmailMap, 0)

	for groupName := range report.GroupsOfDomains {
		groupAuthors, groupLineChanges, yearlyGroupAuthors, yearlyGroupLineChanges := report.accumulateGroupCounts(groupName)
		totalGroupAuthors += groupAuthors
		totalGroupChanges = common.AddLineChanges(totalGroupChanges, groupLineChanges)
		totalGroupYearlyLineChanges.AddYearlyLineChangeMap(yearlyGroupLineChanges)
		totalGroupYearlyAuthors.AddYearlyEmailMap(yearlyGroupAuthors)
	}

	unknownGroupTotalAuthors := report.TotalAuthors - totalGroupAuthors
	unknownGroupTotalLineChanges, _ := common.SubtractLineChanges(report.TotalChanges, totalGroupChanges)
	unknownGroupTotalYearlyLineChanges := report.TotalYearlyLineChanges
	unknownGroupTotalYearlyLineChanges.SubtractYearlyLineChangeMap(totalGroupYearlyLineChanges)
	unknownGroupTotalYearlyAuthors := report.TotalYearlyAuthors
	unknownGroupTotalYearlyAuthors.SubtractYearlyEmailMap(totalGroupYearlyAuthors)

	return NewGroupData(report, fallbackGroupName, unknownGroupTotalAuthors, unknownGroupTotalLineChanges, unknownGroupTotalYearlyLineChanges, unknownGroupTotalYearlyAuthors)
}

func (report *DomainGroupsReport) GroupData(groupName string) *GroupData {
	if groupName == "" || groupName == fallbackGroupName {
		return report.UnknownGroupData()
	}

	totalGroupAuthors, totalGroupLineChanges, totalYearlyGroupAuthors, totalGroupYearlyLineChanges := report.accumulateGroupCounts(groupName)
	return NewGroupData(report, groupName, totalGroupAuthors, totalGroupLineChanges, totalGroupYearlyLineChanges, totalYearlyGroupAuthors)
}
