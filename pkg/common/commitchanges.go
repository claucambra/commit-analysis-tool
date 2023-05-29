package common

type CommitChanges struct {
	Insertions   int
	Deletions    int
	FilesChanged int
}

type YearlyChangeMap map[int]CommitChanges
