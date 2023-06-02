package common

import (
	"math/rand"
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
	setAEmails := generateRandomEmails(generateTestEmailCount)
	setBEmails := generateRandomEmails(generateTestEmailCount)

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
