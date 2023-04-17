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

	authors               map[string]bool
	domainCounts          map[string]int
	corporateEmailDomains map[string]bool
}

func NewCorpAuthorsReport(corporateEmailDomains map[string]bool) *CorpAuthorsReport {
	return &CorpAuthorsReport{
		Commits:            make([]*git.CommitData, 0),
		TotalAuthors:       0,
		NumCorpAuthors:     0,
		CorpAuthorsPercent: 0,
		DomainCountMap:     make(map[string]int),

		authors:               make(map[string]bool),
		domainCounts:          make(map[string]int),
		corporateEmailDomains: corporateEmailDomains,
	}
}

func (report *CorpAuthorsReport) ParseCommits(commits []*git.CommitData) {
	if len(commits) == 0 {
		return
	}

	for _, commit := range commits {
		report.AddCommit(*commit)
	}
}

func (report *CorpAuthorsReport) AddCommit(commit git.CommitData) {
	authorString := commit.AuthorEmail
	if authorString == "" {
		authorString = commit.AuthorName
	}

	if report.authors[authorString] { // Already counted, skip
		return
	} else if authorString != "" {
		report.authors[authorString] = true
		report.TotalAuthors += 1
	}

	splitAuthorEmail := strings.Split(commit.AuthorEmail, "@")

	if len(splitAuthorEmail) < 2 {
		report.domainCounts["unknown"] += 1
		return
	}

	emailDomain := splitAuthorEmail[1]
	report.DomainCountMap[emailDomain] += 1

	if report.corporateEmailDomains[emailDomain] {
		report.NumCorpAuthors += 1
	}

	report.CorpAuthorsPercent = (float32(report.NumCorpAuthors) / float32(report.TotalAuthors)) * 100
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
