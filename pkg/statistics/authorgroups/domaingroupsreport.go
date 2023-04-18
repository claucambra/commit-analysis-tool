package authorgroups

import (
	"strings"

	"github.com/claucambra/commit-analysis-tool/internal/db"
)

const fallbackGroupName = "unknown"

type DomainGroupsReport struct {
	TotalAuthors          int
	TotalCommits          int
	DomainGroups          map[string][]string
	DomainNumAuthors      map[string]int
	DomainNumInsertions   map[string]int
	DomainNumDeletions    map[string]int
	DomainNumFilesChanged map[string]int

	domainToGroup map[string]string
}

func NewDomainGroupsReport(domainGroups map[string][]string) *DomainGroupsReport {
	report := &DomainGroupsReport{
		TotalAuthors:          0,
		TotalCommits:          0,
		DomainGroups:          domainGroups,
		DomainNumAuthors:      map[string]int{},
		DomainNumInsertions:   map[string]int{},
		DomainNumDeletions:    map[string]int{},
		DomainNumFilesChanged: map[string]int{},
		domainToGroup:         map[string]string{},
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

	domainInsertions, domainDeletions, domainFilesChanged, err := db.DomainChanges(authorDomain)
	if err != nil {
		return
	}

	report.DomainNumInsertions[authorDomain] += domainInsertions
	report.DomainNumDeletions[authorDomain] += domainDeletions
	report.DomainNumFilesChanged[authorDomain] += domainFilesChanged
}

func (report *DomainGroupsReport) updateAuthors(authors []string, db *db.SQLiteBackend) {
	for _, author := range authors {
		report.TotalCommits += 1

		if author == "" {
			continue
		}

		splitAuthorEmail := strings.Split(author, "@")
		if len(splitAuthorEmail) < 2 {
			continue
		}

		authorDomain := splitAuthorEmail[1]
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

func (report *DomainGroupsReport) PercentageGroupAuthors(group string) float32 {
	if group == "" {
		return 0
	}

	totalGroupAuthors := 0

	for _, domain := range report.DomainGroups[group] {
		totalGroupAuthors += report.DomainNumAuthors[domain]
	}

	return (float32(totalGroupAuthors) / float32(report.TotalAuthors)) * 100
}
