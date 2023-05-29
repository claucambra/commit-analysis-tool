package authorgroups

import (
	"strings"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

const fallbackDomain = "unknown-domain"
const fallbackGroupName = "unknown"

type DomainGroupsReport struct {
	TotalAuthors int
	TotalCommits int
	TotalChanges *common.LineChanges

	GroupsOfDomains map[string][]string

	DomainTotalNumAuthors   map[string]int
	DomainTotalLineChanges  map[string]*common.LineChanges
	DomainTotalNumDeletions map[string]int
}

func NewDomainGroupsReport(domainGroups map[string][]string) *DomainGroupsReport {
	report := &DomainGroupsReport{
		TotalChanges:                 &common.LineChanges{},
		GroupsOfDomains:              domainGroups,
		DomainTotalNumAuthors:        map[string]int{},
		DomainTotalLineChanges:       map[string]*common.LineChanges{},
		DomainYearlyTotalLineChanges: map[string]common.YearlyLineChangeMap{},
	}

	return report
}

func (report *DomainGroupsReport) updateDomainChanges(authorDomain string, sqlb *db.SQLiteBackend) {
	if authorDomain == "" {
		return
	}

	changes, err := domainChanges(sqlb, authorDomain)
	if err != nil {
		return
	}

	report.TotalChanges.AddLineChanges(&changes.LineChanges)

	if existingDomainLineChanges, ok := report.DomainTotalLineChanges[authorDomain]; ok {
		existingDomainLineChanges.AddLineChanges(&changes.LineChanges)
		report.DomainTotalLineChanges[authorDomain] = existingDomainLineChanges
	} else {
		report.DomainTotalLineChanges[authorDomain] = &changes.LineChanges
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

		report.updateDomainChanges(authorDomain, db)
	}
}

func (report *DomainGroupsReport) Generate(db *db.SQLiteBackend) {
	authors, err := db.Authors()
	if err != nil {
		return
	}

	report.updateAuthors(authors, db)
}

// Returns authors, insertions, deletions
func (report *DomainGroupsReport) accumulateGroupCounts(groupName string) (int, *common.LineChanges) {
	totalGroupAuthors := 0
	totalGroupLineChanges := common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}

	for _, domain := range report.GroupsOfDomains[groupName] {
		reportChanges := report.DomainTotalLineChanges[domain]
		totalGroupLineChanges.AddLineChanges(reportChanges)
		totalGroupAuthors += report.DomainTotalNumAuthors[domain]
	}

	return totalGroupAuthors, &totalGroupLineChanges
}

func (report *DomainGroupsReport) unknownGroupData() *GroupData {
	totalGroupAuthors := 0
	totalGroupChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}

	for groupName := range report.GroupsOfDomains {
		groupAuthors, groupLineChanges := report.accumulateGroupCounts(groupName)
		totalGroupAuthors += groupAuthors
		totalGroupChanges.AddLineChanges(groupLineChanges)
	}

	unknownGroupTotalAuthors := report.TotalAuthors - totalGroupAuthors
	unknownGroupTotalLineChanges := report.TotalChanges
	unknownGroupTotalLineChanges.SubtractLineChanges(totalGroupChanges)

	return NewGroupData(report, fallbackGroupName, unknownGroupTotalAuthors, unknownGroupTotalLineChanges)
}

func (report *DomainGroupsReport) GroupData(groupName string) *GroupData {
	if groupName == "" || groupName == fallbackGroupName {
		return report.unknownGroupData()
	}

	totalGroupAuthors, totalGroupChanges := report.accumulateGroupCounts(groupName)
	return NewGroupData(report, groupName, totalGroupAuthors, totalGroupChanges)
}
