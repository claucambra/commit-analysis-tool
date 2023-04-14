package git

import (
	"fmt"
)

const PrettyFormatSeparator = "__SEPARATOR__"

var PrettyFormat = fmt.Sprintf("%%H%s%%aD%s%%aN%s%%aE%s%%cD%s%%cN%s%%cE",
	PrettyFormatSeparator,
	PrettyFormatSeparator,
	PrettyFormatSeparator,
	PrettyFormatSeparator,
	PrettyFormatSeparator,
	PrettyFormatSeparator)

