package git

import (
	"fmt"
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

