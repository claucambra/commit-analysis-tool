package git

import (
	"strings"
	"time"
)

type CommitData struct {
	id             string
	repoName       string
	authorName     string
	authorEmail    string
	authorTime     time.Time
	committerName  string
	committerEmail string
	committerTime  time.Time
	numInsertions  int
	numDeletions   int
}
