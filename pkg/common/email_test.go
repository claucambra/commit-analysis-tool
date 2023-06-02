package common

import (
	"math/rand"
	"reflect"
	"testing"
)

const generateTestEmailCount = 10

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateRandomEmails(amount int) []string {
	emailArr := make([]string, amount)

	for i := 0; i < amount; i++ {
		email := randSeq(8) + "@" + randSeq(5) + "." + randSeq(3)
		emailArr[i] = email
	}

	return emailArr
}

func generateRandomEmailSets() (EmailSet, EmailSet) {
	commonEmail := "developer@claudiocambra.com"
	setAEmails := append(generateRandomEmails(generateTestEmailCount), commonEmail)
	setBEmails := append(generateRandomEmails(generateTestEmailCount), commonEmail)

	emailSetA := EmailSet{}
	emailSetB := EmailSet{}

	for _, email := range setAEmails {
		emailSetA[email] = true
	}

	for _, email := range setBEmails {
		emailSetB[email] = true
	}

	return emailSetA, emailSetB
}

// EmailSet
func TestAddEmailSet(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	summedEmailSets := AddEmailSet(emailSetA, emailSetB)
	testEmailSet := emailSetA

	for email := range emailSetB {
		testEmailSet[email] = true
	}

	if !reflect.DeepEqual(testEmailSet, summedEmailSets) {
		t.Fatalf(`Added email sets do not match expected email set: 
			Expected %+v
			Received %+v`, testEmailSet, summedEmailSets)
	}
}

func TestSubtractEmailSets(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	subbedEmailSets, _ := SubtractEmailSet(emailSetA, emailSetB)
	testEmailSet := emailSetA

	for email := range emailSetB {
		delete(testEmailSet, email)
	}

	if !reflect.DeepEqual(testEmailSet, subbedEmailSets) {
		t.Fatalf(`Subtracted email sets do not match expected email set: 
			Expected %+v
			Received %+v`, testEmailSet, subbedEmailSets)
	}
}

// YearlyEmailMap
func TestAddEmailSetToYearlyEmailMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYear := 2023
	yem := make(YearlyEmailMap, 0)

	yem.AddEmailSet(emailSetA, testYear)
	if yemAEmailSetA := yem[testYear]; !reflect.DeepEqual(yemAEmailSetA, emailSetA) {
		t.Fatalf(`Added email set to yearly emails map when year not already in map does not match expected changes:
			Expected %+v
			Received %+v`, emailSetA, yemAEmailSetA)
	}

	yem.AddEmailSet(emailSetB, testYear)
	summedEmailSets := AddEmailSet(emailSetA, emailSetB)

	if yemASummedEmailSets := yem[testYear]; !reflect.DeepEqual(yemASummedEmailSets, summedEmailSets) {
		t.Fatalf(`Added email set to yearly emails map when year already in map does not match expected changes:
			Expected %+v
			Received %+v`, summedEmailSets, yemASummedEmailSets)
	}
}

func TestSubtractEmailSetInYearlyEmailsMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYearA := 2023
	testYearB := 2003
	yem := YearlyEmailMap{testYearA: emailSetA}

	yem.SubtractEmailSet(emailSetB, testYearB)
	if _, ok := yem[testYearB]; ok {
		t.Fatalf("Subtracting email set from a year not present in YEM should not add this year to YEM.")
	}

	yem = YearlyEmailMap{testYearA: emailSetA}
	yem.SubtractEmailSet(emailSetB, testYearA)

	expectedSubEmailSet, _ := SubtractEmailSet(emailSetA, emailSetB)

	if subEmailSet := yem[testYearA]; !reflect.DeepEqual(subEmailSet, expectedSubEmailSet) {
		t.Fatalf(`Subtracted email set from yearly email map does not match expected changes:
			Expected %+v
			Received %+v`, expectedSubEmailSet, subEmailSet)
	}
}

func TestAddYearlyEmailMapToYearlyEmailMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYearA := 2023
	testYearB := 2003
	yemA := YearlyEmailMap{testYearA: emailSetA}
	yemB := YearlyEmailMap{testYearB: emailSetB}

	yemA.AddYearlyEmailMap(yemB)

	if addedEmailSetB, ok := yemA[testYearB]; !ok {
		t.Fatalf("Adding email set from a year of YEM B not present in YEM A should add this year to YEM A.")
	} else if !reflect.DeepEqual(addedEmailSetB, emailSetB) {
		t.Fatalf(`Added YLCM B to YLCM A does not match expected line changes:
			Expected %+v
			Received %+v`, emailSetB, addedEmailSetB)
	}

	yemB = YearlyEmailMap{testYearA: emailSetB}
	yemA.AddYearlyEmailMap(yemB)

	expectedSummedEmailSets := AddEmailSet(emailSetA, emailSetB)

	if addedEmailSets := yemA[testYearA]; !reflect.DeepEqual(addedEmailSets, expectedSummedEmailSets) {
		t.Fatalf(`Added YLCM B email set to YLCM A does not match expected email set:
			Expected %+v
			Received %+v`, addedEmailSets, expectedSummedEmailSets)
	}
}

func TestSubtractYearlyEmailMapToYearlyEmailMap(t *testing.T) {
	emailSetA, emailSetB := generateRandomEmailSets()
	testYearA := 2023
	testYearB := 2003
	yemA := YearlyEmailMap{testYearA: emailSetA}
	yemB := YearlyEmailMap{testYearB: emailSetB}

	yemA.SubtractYearlyEmailMap(yemB)

	if _, ok := yemA[testYearB]; ok {
		t.Fatalf("Subtracting YEM B email set from a year not present in YEM A should not add this year to YEM A.")
	}

	yemA.AddEmailSet(emailSetA, testYearB)
	yemA.SubtractYearlyEmailMap(yemB)

	expectedSubbedEmailSets, _ := SubtractEmailSet(emailSetA, emailSetB)

	if testYearBSubEmailSet, ok := yemA[testYearB]; !ok {
		t.Fatalf("Test year B should now be present in YLCM A")
	} else if !reflect.DeepEqual(testYearBSubEmailSet, expectedSubbedEmailSets) {
		t.Fatalf(`Subtracted YLCM B email set to yearly email set map does not match expected changes:
			Expected %+v
			Received %+v`, expectedSubbedEmailSets, testYearBSubEmailSet)
	}
}
