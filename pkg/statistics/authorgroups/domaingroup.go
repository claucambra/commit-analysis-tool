package authorgroups

import (
	"fmt"
	"sort"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

type DomainGroup struct {
	Name          string
	AuthorCount   int
	DomainAuthors map[string][]*common.Author
}

func (group *DomainGroup) domainsString() string {
	// Get sorted domains
	sortedDomainNames := make([]string, 0, len(group.DomainAuthors))
	for domainName := range group.DomainAuthors {
		sortedDomainNames = append(sortedDomainNames, domainName)
	}

	sort.SliceStable(sortedDomainNames, func(i, j int) bool {
		domainA := sortedDomainNames[i]
		domainB := sortedDomainNames[j]

		domainACount := len(group.DomainAuthors[domainA])
		domainBCount := len(group.DomainAuthors[domainB])

		if domainACount == domainBCount {
			return domainA < domainB
		}

		return domainACount > domainBCount
	})

	reportString := ""

	for _, domainName := range sortedDomainNames {
		print(domainName)
		domainCount := len(group.DomainAuthors[domainName])
		reportString += fmt.Sprintf("\t\t%s:\t%d\n", domainName, domainCount)
	}

	return reportString
}
