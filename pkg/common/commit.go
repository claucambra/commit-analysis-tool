package common

type Commit struct {
	Changes

	Id            string
	RepoName      string
	Author        Person
	AuthorTime    int64
	Committer     Person
	CommitterTime int64
}

// Similar to RFC1123Z but without trailing zero on day
const TimeFormat = "Mon, 2 Jan 2006 15:04:05 -0700"
