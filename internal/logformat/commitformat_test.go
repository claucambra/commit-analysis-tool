package logformat

import (
	"testing"
)

func TestPrettyFormat(t *testing.T) {
	expectedPrettyFormat := "PRETTYFORMATSTART__%H__SEPARATOR__%aD__SEPARATOR__%aN__SEPARATOR__%aE__SEPARATOR__%cD__SEPARATOR__%cN__SEPARATOR__%cE__SEPARATOR__%s__SEPARATOR__%b__PRETTYFORMATEND"
	prettyFormat := PrettyFormatString()
	if prettyFormat != expectedPrettyFormat {
		t.Fatalf(`Received incorrect PrettyFormat.
			Expected: %s
			Received: %s`, expectedPrettyFormat, prettyFormat)
	}
}

func TestPrettyFormatParameterCount(t *testing.T) {
	expectedParameterCount := 9
	parameterCount := PrettyFormatStringParameterCount()
	if parameterCount != expectedParameterCount {
		t.Fatalf("Received incorrect number of parameters for PrettyFormatParameterCount: %d, expected %d", parameterCount, expectedParameterCount)
	}
}
