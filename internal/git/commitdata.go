package git

type CommitData struct {
	id              string
	repoName        string
	authorName      string
	authorEmail     string
	authorTime      int64
	committerName   string
	committerEmail  string
	committerTime   int64
	numInsertions   int
	numDeletions    int
	numFilesChanged int
}

// Similar to RFC1123Z but without trailing zero on day
const TimeFormat = "Mon, 2 Jan 2006 15:04:05 -0700"
