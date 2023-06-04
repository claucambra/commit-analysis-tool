package common

// Similar to RFC1123Z but without trailing zero on day
const TimeFormat = "Mon, 2 Jan 2006 15:04:05 -0700"

type Commit struct {
	Changes

	Id            string
	RepoName      string
	Author        Person
	AuthorTime    int64
	Committer     Person
	CommitterTime int64
	Subject       string
	Body          string
}

type CommitMap map[string]*Commit

func (cm *CommitMap) AddCommitMap(cmToAdd CommitMap) {
	for id, commit := range cmToAdd {
		(*cm)[id] = commit
	}
}

func (cm *CommitMap) SubtractCommitMap(cmToSubtract CommitMap) {
	for id := range cmToSubtract {
		delete(*cm, id)
	}
}
