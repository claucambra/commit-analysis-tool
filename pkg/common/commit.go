package common

import (
	"sort"
	"time"
)

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

func (cm *CommitMap) YearRange(excludeEmpty bool) []int {
	years := []int{}

	for _, commit := range *cm {
		commitTime := time.Unix(commit.AuthorTime, 0).UTC()
		commitYear := commitTime.Year()

		if found, _ := SliceContains(years, commitYear); found {
			continue
		}

		years = append(years, commitYear)
	}

	sort.Slice(years, func(i, j int) bool {
		return years[i] < years[j]
	})

	if excludeEmpty || len(years) == 0 {
		return years
	}

	firstYear := years[0]
	lastYear := years[len(years)-1]

	filledYears := []int{}
	for i := firstYear; i <= lastYear; i++ {
		filledYears = append(filledYears, i)
	}

	return filledYears
}
