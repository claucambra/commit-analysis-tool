package logformat

import (
	"testing"
)

func TestPrettyFormat(t *testing.T) {
	expectedPrettyFormat := "%H__SEPARATOR__%aD__SEPARATOR__%aN__SEPARATOR__%aE__SEPARATOR__%cD__SEPARATOR__%cN__SEPARATOR__%cE"
	prettyFormat := PrettyFormatString()
	if prettyFormat != expectedPrettyFormat {
		t.Fatalf(`Received incorrect PrettyFormat.
			Expected: %s
			Received: %s`, expectedPrettyFormat, prettyFormat)
	}
}

func TestPrettyFormatParameterCount(t *testing.T) {
	expectedParameterCount := 7
	parameterCount := PrettyFormatStringParameterCount()
	if parameterCount != expectedParameterCount {
		t.Fatalf("Received incorrect number of parameters for PrettyFormatParameterCount: %d, expected %d", parameterCount, expectedParameterCount)
	}
}
