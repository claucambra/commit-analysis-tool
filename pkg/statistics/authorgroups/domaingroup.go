package authorgroups

import (
	"fmt"
	"sort"
)

type DomainGroup struct {
	AuthorCount  int
	DomainCounts map[string]int
}

func (group *DomainGroup) domainsString() string {
	// Get sorted domains
	sortedDomainNames := make([]string, 0, len(group.DomainCounts))
	for domainName := range group.DomainCounts {
		sortedDomainNames = append(sortedDomainNames, domainName)
	}

	sort.SliceStable(sortedDomainNames, func(i, j int) bool {
		domainA := sortedDomainNames[i]
		domainB := sortedDomainNames[j]

		domainACount := group.DomainCounts[domainA]
		domainBCount := group.DomainCounts[domainB]

		if domainACount == domainBCount {
			return domainA < domainB
		}

		return domainACount > domainBCount
	})

	reportString := ""

	for _, domainName := range sortedDomainNames {
		print(domainName)
		DomainCounts := group.DomainCounts[domainName]
		reportString += fmt.Sprintf("\t\t%s:\t%d\n", domainName, DomainCounts)
	}

	return reportString
}
