package statistics

import (
	"fmt"
	"strings"

	"github.com/claucambra/commit-analysis-tool/internal/git"
)

type CorpAuthorsReport struct {
	Commits            []*git.CommitData
	TotalAuthors       int
	NumCorpAuthors     int
	CorpAuthorsPercent float32
	DomainCountMap     map[string]int
}

func NewCorpAuthorsReport(commits []*git.CommitData, corporateEmailDomains map[string]bool) *CorpAuthorsReport {
	report := &CorpAuthorsReport{
		Commits:            commits,
		TotalAuthors:       0,
		NumCorpAuthors:     0,
		CorpAuthorsPercent: 0,
		DomainCountMap:     make(map[string]int),
	}

	if len(commits) == 0 {
		return report
	}

	authorsSet := make(map[string]bool)
	domainCounts := make(map[string]int)

	for _, commit := range commits {
		authorString := commit.AuthorEmail
		if authorString == "" {
			authorString = commit.AuthorName
		}

		if authorsSet[authorString] { // Already counted, skip
			continue
		} else if authorString != "" {
			authorsSet[authorString] = true
			report.TotalAuthors += 1
		}

		splitAuthorEmail := strings.Split(commit.AuthorEmail, "@")

		if len(splitAuthorEmail) < 2 {
			domainCounts["unknown"] += 1
			continue
		}

		emailDomain := splitAuthorEmail[1]
		report.DomainCountMap[emailDomain] += 1

		if corporateEmailDomains[emailDomain] {
			report.NumCorpAuthors += 1
		}
	}

	report.CorpAuthorsPercent = (float32(report.NumCorpAuthors) / float32(report.TotalAuthors)) * 100
	return report
}

func (report *CorpAuthorsReport) String() string {
	reportString := "Corporate authors report\n"
	reportString += fmt.Sprintf("Total repository authors: %d\n", report.TotalAuthors)
	reportString += fmt.Sprintf("Number of corporate authors: %d (%f%%)\n", report.NumCorpAuthors, report.CorpAuthorsPercent)
	reportString += "Number of authors by domain:\n"

	for domain, count := range report.DomainCountMap {
		reportString += fmt.Sprintf("\t%s: %d\n", domain, count)
	}

	return reportString
}
