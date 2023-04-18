package authorgroups

import "github.com/claucambra/commit-analysis-tool/pkg/common"

type DomainGroup struct {
	AuthorCount   int
	DomainAuthors map[string][]*common.Author
}
