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
	TotalAuthors common.EmailSet
	TotalChanges *common.LineChanges
	TotalCommits common.CommitMap

	GroupsOfDomains map[string][]string

	DomainTotalAuthors     map[string]common.EmailSet
	DomainTotalLineChanges map[string]*common.LineChanges

	DomainCommits map[string]common.CommitMap

	sqlb *db.SQLiteBackend
}

func NewDomainGroupsReport(domainGroups map[string][]string, sqlb *db.SQLiteBackend) *DomainGroupsReport {
	return &DomainGroupsReport{
		TotalAuthors:           common.EmailSet{},
		TotalChanges:           &common.LineChanges{},
		TotalCommits:           common.CommitMap{},
		GroupsOfDomains:        domainGroups,
		DomainTotalAuthors:     map[string]common.EmailSet{},
		DomainTotalLineChanges: map[string]*common.LineChanges{},
		DomainCommits:          map[string]common.CommitMap{},
		sqlb:                   sqlb,
	}
}

func (report *DomainGroupsReport) resetStats() {
	report.TotalAuthors = common.EmailSet{}
	report.TotalChanges = &common.LineChanges{}
	report.TotalCommits = common.CommitMap{}
	report.DomainTotalAuthors = map[string]common.EmailSet{}
	report.DomainTotalLineChanges = map[string]*common.LineChanges{}
	report.DomainCommits = map[string]common.CommitMap{}
}

func (report *DomainGroupsReport) updateDomainChanges() {
	for authorDomain := range report.DomainTotalAuthors {
		log.Printf("Updating domain groups report data for domain: %s", authorDomain)

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

		domainCommits, err := domainCommits(report.sqlb, authorDomain)
		if err != nil {
			log.Fatalf("Error retrieving commits for domain %s, received error: %s", authorDomain, err)
			return
		}

		for _, commit := range domainCommits {
			report.TotalCommits[commit.Id] = commit

			if _, ok := report.DomainCommits[authorDomain]; !ok {
				report.DomainCommits[authorDomain] = common.CommitMap{commit.Id: commit}
			} else {
				report.DomainCommits[authorDomain][commit.Id] = commit
			}
		}
	}
}

func (report *DomainGroupsReport) updateAuthors(authors []string) {
	log.Printf("Updating domain groups report authors.")

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

	log.Println("Generating domain groups report.")

	report.resetStats()
	report.updateAuthors(authors)
	report.updateDomainChanges()
}

// Returns authors, insertions, deletions
func (report *DomainGroupsReport) accumulateGroupCounts(groupName string) (
	common.EmailSet,
	*common.LineChanges,
	common.CommitMap) {

	totalGroupAuthors := common.EmailSet{}
	totalGroupLineChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}
	totalGroupCommits := common.CommitMap{}

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

		totalGroupCommits.AddCommitMap(report.DomainCommits[domain])
	}

	return totalGroupAuthors,
		totalGroupLineChanges,
		totalGroupCommits
}

func (report *DomainGroupsReport) UnknownGroupData() *GroupData {
	totalGroupAuthors := common.EmailSet{}
	totalGroupChanges := &common.LineChanges{
		NumInsertions: 0,
		NumDeletions:  0,
	}
	totalGroupCommits := common.CommitMap{}

	for groupName := range report.GroupsOfDomains {
		groupAuthors, groupLineChanges, groupCommits := report.accumulateGroupCounts(groupName)
		totalGroupAuthors = common.AddEmailSet(totalGroupAuthors, groupAuthors)
		totalGroupChanges = common.AddLineChanges(totalGroupChanges, groupLineChanges)
		totalGroupCommits.AddCommitMap(groupCommits)
	}

	unknownGroupTotalAuthors, _ := common.SubtractEmailSet(report.TotalAuthors, totalGroupAuthors)
	unknownGroupTotalLineChanges, _ := common.SubtractLineChanges(report.TotalChanges, totalGroupChanges)

	var unknownGroupCommits common.CommitMap = common.CopyMap(report.TotalCommits)
	unknownGroupCommits.SubtractCommitMap(totalGroupCommits)

	return NewGroupData(report,
		fallbackGroupName,
		unknownGroupTotalAuthors,
		unknownGroupTotalLineChanges,
		unknownGroupCommits)
}

func (report *DomainGroupsReport) GroupData(groupName string) *GroupData {
	if groupName == "" || groupName == fallbackGroupName {
		return report.UnknownGroupData()
	}

	totalGroupAuthors, totalGroupLineChanges, totalGroupCommits := report.accumulateGroupCounts(groupName)

	return NewGroupData(report,
		groupName,
		totalGroupAuthors,
		totalGroupLineChanges,
		totalGroupCommits)
}
