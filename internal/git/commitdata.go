package git

type CommitData struct {
	id             string
	repoName       string
	authorName     string
	authorEmail    string
	authorTime     int64
	committerName  string
	committerEmail string
	committerTime  int64
	numInsertions  int
	numDeletions   int
}

