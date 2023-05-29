package common

type Report interface {
	AddCommit(Commit)
	String() string
}
