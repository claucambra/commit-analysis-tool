package common

type Person struct {
	Name  string
	Email string
}

type EmailSet map[string]bool
type YearlyEmailMap map[int]EmailSet

func AddEmailSet(a EmailSet, b EmailSet) EmailSet {
	newSet := a

	for email := range b {
		newSet[email] = true
	}

	return newSet
}

func SubtractEmailSet(a EmailSet, b EmailSet) EmailSet {
	newSet := a

	for email := range b {
		delete(newSet, email)
	}

	return newSet
}

func (yem *YearlyEmailMap) CountArray(years []int) []int {
	yearsToReturn := years

	if yearsToReturn == nil {
		yearsToReturn = SortedMapKeys(*yem)
	}

	countArray := make([]int, len(yearsToReturn))

	for i, year := range yearsToReturn {
		countArray[i] = len((*yem)[year])
	}

	return countArray
}

func (yem *YearlyEmailMap) AddEmailSet(emailSetToAdd EmailSet, year int) {
	AdditiveValueMapInsert[int, EmailSet, YearlyEmailMap](*yem, year, AddEmailSet, emailSetToAdd)
}

func (yem *YearlyEmailMap) SubtractEmailSet(emailSetToAdd EmailSet, year int) {
	AdditiveValueMapInsert[int, EmailSet, YearlyEmailMap](*yem, year, AddEmailSet, emailSetToAdd)
}

func (yem *YearlyEmailMap) AddYearlyPeopleMap(yemToAdd YearlyEmailMap) {
	for year, emailsToAdd := range yemToAdd {
		if existingEmails, ok := (*yem)[year]; ok {
			(*yem)[year] = AddEmailSet(existingEmails, emailsToAdd)
		} else {
			(*yem)[year] = emailsToAdd
		}
	}
}

func (yem *YearlyEmailMap) SubtractYearlyPeopleMap(yemToSubtract YearlyEmailMap) {
	for year, emailsToSubtract := range yemToSubtract {
		if existingEmails, ok := (*yem)[year]; ok {
			subtractedEmails := SubtractEmailSet(existingEmails, emailsToSubtract)
			(*yem)[year] = subtractedEmails
		}
	}
}
