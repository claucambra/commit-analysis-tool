package authorgroups

type GroupData struct {
	GroupName string

	NumAuthors    int
	NumInsertions int
	NumDeletions  int

	AuthorsPercent    float32
	InsertionsPercent float32
	DeletionsPercent  float32
}
