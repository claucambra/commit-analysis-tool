package authorgroups

import (
	"log"
	"regexp"
	"strings"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

const fallbackDomain = "unknown-domain"
const fallbackGroupName = "unknown"

// Report of the organised raw data around a grouping of domains
type DomainGroupsReport struct {
	TotalAuthors           common.EmailSet
	TotalChanges           *common.LineChanges
	TotalYearlyAuthors     common.YearlyEmailMap
	TotalYearlyLineChanges common.YearlyLineChangeMap

	GroupsOfDomains map[string][]string

	DomainTotalAuthors      map[string]common.EmailSet
	DomainTotalLineChanges  map[string]*common.LineChanges
	DomainTotalNumDeletions map[string]int

	DomainTotalYearlyAuthors     map[string]common.YearlyEmailMap
	DomainTotalYearlyLineChanges map[string]common.YearlyLineChangeMap

	sqlb *db.SQLiteBackend
}

func NewDomainGroupsReport(domainGroups map[string][]string, sqlb *db.SQLiteBackend) *DomainGroupsReport {
	return &DomainGroupsReport{
		TotalAuthors:                 common.EmailSet{},
		TotalChanges:                 &common.LineChanges{},
		TotalYearlyAuthors:           common.YearlyEmailMap{},
		TotalYearlyLineChanges:       common.YearlyLineChangeMap{},
		GroupsOfDomains:              domainGroups,
		DomainTotalAuthors:           map[string]common.EmailSet{},
		DomainTotalLineChanges:       map[string]*common.LineChanges{},
		DomainTotalYearlyAuthors:     map[string]common.YearlyEmailMap{},
		DomainTotalYearlyLineChanges: map[string]common.YearlyLineChangeMap{},
		sqlb:                         sqlb,
	}
}

func (report *DomainGroupsReport) resetStats() {
	report.TotalAuthors = common.EmailSet{}
	report.TotalChanges = &common.LineChanges{}
	report.TotalYearlyLineChanges = common.YearlyLineChangeMap{}
	report.DomainTotalAuthors = map[string]common.EmailSet{}
	report.DomainTotalLineChanges = map[string]*common.LineChanges{}
	report.DomainTotalYearlyLineChanges = map[string]common.YearlyLineChangeMap{}
}

func (report *DomainGroupsReport) updateDomainChanges() {
	for authorDomain := range report.DomainTotalAuthors {
		lineChanges, err := domainLineChanges(report.sqlb, authorDomain)
		if err != nil {
			log.Fatalf("Error retrieving line changes for domain %s, received error: %s", authorDomain, err)
			return
		}

		report.TotalChanges = common.AddLineChanges(report.TotalChanges, lineChanges)

		if existingDomainLineChanges, ok := report.DomainTotalLineChanges[authorDomain]; ok {
			summedDomainLineChanges := common.AddLineChanges(existingDomainLineChanges, lineChanges)
			report.DomainTotalLineChanges[authorDomain] = summedDomainLineChanges
		} else {
			report.DomainTotalLineChanges[authorDomain] = lineChanges
		}

		yearlyLineChanges, err := domainYearlyLineChanges(report.sqlb, authorDomain)
		if err != nil {
			log.Fatalf("Error retrieving yearly line changes for domain %s, received error: %s", authorDomain, err)
			return
		}

		report.TotalYearlyLineChanges.AddYearlyLineChangeMap(yearlyLineChanges)

		if existingDomainYearLineChanges, ok := report.DomainTotalYearlyLineChanges[authorDomain]; ok {
			existingDomainYearLineChanges.AddYearlyLineChangeMap(yearlyLineChanges)
			report.DomainTotalYearlyLineChanges[authorDomain] = existingDomainYearLineChanges
		} else {
			report.DomainTotalYearlyLineChanges[authorDomain] = yearlyLineChanges
		}

		yearlyAuthors, err := domainYearlyAuthors(report.sqlb, authorDomain)
		if err != nil {
			log.Fatalf("Error retrieving yearly authors for domain %s, received error: %s", authorDomain, err)
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

func (report *DomainGroupsReport) updateAuthors(authors []string) {
	for _, author := range authors {
		if author == "" {
			continue
		}

		authorDomain := fallbackDomain
		splitAuthorEmail := strings.Split(author, "@")

		if len(splitAuthorEmail) >= 2 {
			authorDomain = splitAuthorEmail[1]
		}

		currentDomainAuthors := report.DomainTotalAuthors[authorDomain]
		report.DomainTotalAuthors[authorDomain] = common.AddEmailSet(currentDomainAuthors, common.EmailSet{author: true})
		report.TotalAuthors[author] = true
	}
}

func (report *DomainGroupsReport) Generate() {
	authors, err := report.sqlb.Authors()
	if err != nil {
		return
	}

	report.resetStats()
	report.updateAuthors(authors)
	report.updateDomainChanges()
}

// Returns authors, insertions, deletions
func (report *DomainGroupsReport) accumulateGroupCounts(groupName string) (common.EmailSet, *common.LineChanges, common.YearlyEmailMap, common.YearlyLineChangeMap) {
	totalGroupAuthors := common.EmailSet{}
	totalGroupLineChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}
	totalGroupYearlyLineChanges := make(common.YearlyLineChangeMap, 0)
	totalGroupYearlyAuthors := make(common.YearlyEmailMap, 0)

	// Slice of input group of domains string that matches actual domains extracted from emails
	matchingDomains := []string{}

	// Find all domains that match domains in group, treat domains in group as potential regexes
	for _, groupDomainString := range report.GroupsOfDomains[groupName] {
		groupDomainStringRegex := regexp.MustCompile(groupDomainString)

		for domain := range report.DomainTotalAuthors {
			if groupDomainStringRegex.MatchString(domain) {
				matchingDomains = append(matchingDomains, domain)
			}
		}
	}

	for _, domain := range matchingDomains {
		reportChanges, ok := report.DomainTotalLineChanges[domain]
		if !ok {
			continue
		}

		totalGroupLineChanges = common.AddLineChanges(totalGroupLineChanges, reportChanges)
		totalGroupAuthors = common.AddEmailSet(totalGroupAuthors, report.DomainTotalAuthors[domain])

		reportYearlyChanges := report.DomainTotalYearlyLineChanges[domain]
		totalGroupYearlyLineChanges.AddYearlyLineChangeMap(reportYearlyChanges)

		reportYearlyAuthors := report.DomainTotalYearlyAuthors[domain]
		totalGroupYearlyAuthors.AddYearlyEmailMap(reportYearlyAuthors)
	}

	return totalGroupAuthors, totalGroupLineChanges, totalGroupYearlyAuthors, totalGroupYearlyLineChanges
}

func (report *DomainGroupsReport) UnknownGroupData() *GroupData {
	totalGroupAuthors := common.EmailSet{}
	totalGroupChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}
	totalGroupYearlyLineChanges := make(common.YearlyLineChangeMap, 0)
	totalGroupYearlyAuthors := make(common.YearlyEmailMap, 0)

	for groupName := range report.GroupsOfDomains {
		groupAuthors, groupLineChanges, yearlyGroupAuthors, yearlyGroupLineChanges := report.accumulateGroupCounts(groupName)
		totalGroupAuthors = common.AddEmailSet(totalGroupAuthors, groupAuthors)
		totalGroupChanges = common.AddLineChanges(totalGroupChanges, groupLineChanges)
		totalGroupYearlyLineChanges.AddYearlyLineChangeMap(yearlyGroupLineChanges)
		totalGroupYearlyAuthors.AddYearlyEmailMap(yearlyGroupAuthors)
	}

	unknownGroupTotalAuthors, _ := common.SubtractEmailSet(report.TotalAuthors, totalGroupAuthors)
	unknownGroupTotalLineChanges, _ := common.SubtractLineChanges(report.TotalChanges, totalGroupChanges)
	unknownGroupTotalYearlyLineChanges := report.TotalYearlyLineChanges
	unknownGroupTotalYearlyLineChanges.SubtractYearlyLineChangeMap(totalGroupYearlyLineChanges)
	unknownGroupTotalYearlyAuthors := report.TotalYearlyAuthors
	unknownGroupTotalYearlyAuthors.SubtractYearlyEmailMap(totalGroupYearlyAuthors)

	return NewGroupData(report,
		fallbackGroupName,
		unknownGroupTotalAuthors,
		unknownGroupTotalLineChanges,
		unknownGroupTotalYearlyLineChanges,
		unknownGroupTotalYearlyAuthors)
}

func (report *DomainGroupsReport) GroupData(groupName string) *GroupData {
	if groupName == "" || groupName == fallbackGroupName {
		return report.UnknownGroupData()
	}

	totalGroupAuthors, totalGroupLineChanges, totalYearlyGroupAuthors, totalGroupYearlyLineChanges := report.accumulateGroupCounts(groupName)

	return NewGroupData(report,
		groupName,
		totalGroupAuthors,
		totalGroupLineChanges,
		totalGroupYearlyLineChanges,
		totalYearlyGroupAuthors)
}
