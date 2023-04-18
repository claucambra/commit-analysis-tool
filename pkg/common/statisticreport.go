package common

type Report interface {
	AddCommit(CommitData)
	String() string
}
