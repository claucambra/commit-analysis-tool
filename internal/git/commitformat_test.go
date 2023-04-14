package git

import (
	"testing"
)

func TestPrettyFormat(t *testing.T) {
	expectedPrettyFormat := "%H__SEPARATOR__%aD__SEPARATOR__%aN__SEPARATOR__%aE__SEPARATOR__%cD__SEPARATOR__%cN__SEPARATOR__%cE"
	if PrettyFormat != expectedPrettyFormat {
		t.Fatalf(`Received incorrect PrettyFormat.
			Expected: %s
			Received: %s`, expectedPrettyFormat, PrettyFormat)
	}
}

