package authorgroups

import (
	"strings"

	"github.com/claucambra/commit-analysis-tool/internal/db"
)

const fallbackDomain = "unknown-domain"
const fallbackGroupName = "unknown"

type DomainGroupsReport struct {
	TotalAuthors    int
	TotalCommits    int
	TotalInsertions int
	TotalDeletions  int

	DomainGroups        map[string][]string
	DomainNumAuthors    map[string]int
	DomainNumInsertions map[string]int
	DomainNumDeletions  map[string]int

	domainToGroup map[string]string
}

func NewDomainGroupsReport(domainGroups map[string][]string) *DomainGroupsReport {
	report := &DomainGroupsReport{
		TotalAuthors:    0,
		TotalCommits:    0,
		TotalInsertions: 0,
		TotalDeletions:  0,

		DomainGroups:        domainGroups,
		DomainNumAuthors:    map[string]int{},
		DomainNumInsertions: map[string]int{},
		DomainNumDeletions:  map[string]int{},

		domainToGroup: map[string]string{},
	}

	for groupName, domainNames := range domainGroups {
		for _, domainName := range domainNames {
			report.domainToGroup[domainName] = groupName
		}
	}

	return report
}

func (report *DomainGroupsReport) updateDomainChanges(authorDomain string, db *db.SQLiteBackend) {
	if authorDomain == "" {
		return
	}

	domainInsertions, domainDeletions, _, err := db.DomainChanges(authorDomain)
	if err != nil {
		return
	}

	report.TotalInsertions += domainInsertions
	report.TotalDeletions += domainDeletions

	report.DomainNumInsertions[authorDomain] += domainInsertions
	report.DomainNumDeletions[authorDomain] += domainDeletions
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

		report.DomainNumAuthors[authorDomain] += 1
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
func (report *DomainGroupsReport) accumulateGroupCounts(groupName string) (int, int, int) {
	totalGroupAuthors := 0
	totalGroupInsertions := 0
	totalGroupDeletions := 0

	for _, domain := range report.DomainGroups[groupName] {
		totalGroupAuthors += report.DomainNumAuthors[domain]
		totalGroupInsertions += report.DomainNumInsertions[domain]
		totalGroupDeletions += report.DomainNumDeletions[domain]
	}

	return totalGroupAuthors, totalGroupInsertions, totalGroupDeletions
}

func (report *DomainGroupsReport) unknownGroupData() *GroupData {
	totalGroupAuthors := 0
	totalGroupInsertions := 0
	totalGroupDeletions := 0

	for groupName := range report.DomainGroups {
		domainAuthors, domainInserts, domainDeletes := report.accumulateGroupCounts(groupName)
		totalGroupAuthors += domainAuthors
		totalGroupInsertions += domainInserts
		totalGroupDeletions += domainDeletes
	}

	unknownGroupTotalAuthors := report.TotalAuthors - totalGroupAuthors
	unknownGroupTotalInsertions := report.TotalInsertions - totalGroupInsertions
	unknownGroupTotalDeletions := report.TotalDeletions - totalGroupDeletions

	groupData := new(GroupData)
	groupData.GroupName = fallbackGroupName
	groupData.NumAuthors = unknownGroupTotalAuthors
	groupData.NumInsertions = unknownGroupTotalInsertions
	groupData.NumDeletions = unknownGroupTotalDeletions
	groupData.AuthorsPercent = (float32(unknownGroupTotalAuthors) / float32(report.TotalAuthors)) * 100
	groupData.InsertionsPercent = (float32(unknownGroupTotalInsertions) / float32(report.TotalInsertions)) * 100
	groupData.DeletionsPercent = (float32(unknownGroupTotalDeletions) / float32(report.TotalDeletions)) * 100

	return groupData
}

func (report *DomainGroupsReport) GroupData(groupName string) *GroupData {
	if groupName == "" || groupName == fallbackGroupName {
		return report.unknownGroupData()
	}

	totalGroupAuthors, totalGroupInsertions, totalGroupDeletions := report.accumulateGroupCounts(groupName)

	groupData := new(GroupData)
	groupData.GroupName = groupName
	groupData.NumAuthors = totalGroupAuthors
	groupData.NumInsertions = totalGroupInsertions
	groupData.NumDeletions = totalGroupDeletions
	groupData.AuthorsPercent = (float32(totalGroupAuthors) / float32(report.TotalAuthors)) * 100
	groupData.InsertionsPercent = (float32(totalGroupInsertions) / float32(report.TotalInsertions)) * 100
	groupData.DeletionsPercent = (float32(totalGroupDeletions) / float32(report.TotalDeletions)) * 100

	return groupData
}
