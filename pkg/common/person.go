package common

type Person struct {
	Name  string
	Email string
}

type YearlyPeopleMap map[int][]*Person
