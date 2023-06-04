package logformat

import (
	"fmt"
	"strings"
)

const PrettyFormatStringStart = "PRETTYFORMATSTART__"
const PrettyFormatStringSeparator = "__SEPARATOR__"
const PrettyFormatStringEnd = "__PRETTYFORMATEND"

func PrettyFormatString() string {
	return fmt.Sprintf("%s%%H%s%%aD%s%%aN%s%%aE%s%%cD%s%%cN%s%%cE%s%%s%s%%b%s",
		PrettyFormatStringStart,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringEnd)
}

func PrettyFormatStringParameterCount() int {
	splitPrettyFormat := strings.Split(PrettyFormatString(), PrettyFormatStringSeparator)
	return len(splitPrettyFormat)
}
