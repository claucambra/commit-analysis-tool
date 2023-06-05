package commitcoding

import (
	"log"
	"regexp"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

type CommitCodingReport struct {
	Commits          common.CommitMap
	CodeMap          map[string][]string // Mapping of string matches to features
	CodeMatchCommits map[string][]*common.Commit
}

func NewCommitCodingReport(commits common.CommitMap, codeMap map[string][]string) *CommitCodingReport {
	return &CommitCodingReport{
		Commits:          commits,
		CodeMap:          codeMap,
		CodeMatchCommits: map[string][]*common.Commit{},
	}
}

func (ccr *CommitCodingReport) Generate() {
	for codeCategory, regexStringSlice := range ccr.CodeMap {
		log.Printf("Finding commit matches for coding analysis category: %s", codeCategory)

		codeCategoryCommits := []*common.Commit{}

		for _, regexString := range regexStringSlice {
			regex := regexp.MustCompile(regexString)

			for _, commit := range ccr.Commits {
				fullCommitBody := commit.Subject + "\n" + commit.Body

				if regex.MatchString(fullCommitBody) {
					codeCategoryCommits = append(codeCategoryCommits, commit)
				}
			}
		}

		ccr.CodeMatchCommits[codeCategory] = codeCategoryCommits
	}
}
