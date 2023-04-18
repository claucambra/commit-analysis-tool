package common

type CommitData struct {
	Id              string
	RepoName        string
	Author          *Author
	AuthorTime      int64
	CommitterName   string
	CommitterEmail  string
	CommitterTime   int64
	NumInsertions   int
	NumDeletions    int
	NumFilesChanged int
}

// Similar to RFC1123Z but without trailing zero on day
const TimeFormat = "Mon, 2 Jan 2006 15:04:05 -0700"
