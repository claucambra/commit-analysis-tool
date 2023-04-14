package git

import (
	"fmt"
	"strings"
)

const PrettyFormatStringSeparator = "__SEPARATOR__"

func PrettyFormatString() string {
	return fmt.Sprintf("%%H%s%%aD%s%%aN%s%%aE%s%%cD%s%%cN%s%%cE",
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator,
		PrettyFormatStringSeparator)
}

func PrettyFormatStringParameterCount() int {
	splitPrettyFormat := strings.Split(PrettyFormatString(), PrettyFormatStringSeparator)
	return len(splitPrettyFormat)
}
