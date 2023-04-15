package git

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var numberRegex = regexp.MustCompile("[0-9]+")
var commitRegex = regexp.MustCompile(fmt.Sprintf("^[0-9a-f]+%s", PrettyFormatStringSeparator))
var insertionsRegex = regexp.MustCompile("([0-9]+) insertions?")
var deletionsRegex = regexp.MustCompile("([0-9]+) deletions?")

/**
 * Convenience function to get a specific number of changes in the changes line (i.e. 320 insertions).
 */
func parseChangesLine(changesLogLine string, specificChangesRegex *regexp.Regexp) (int, error) {
	if !specificChangesRegex.MatchString(changesLogLine) {
		return 0, nil
	}

	specificChangesString := specificChangesRegex.FindString(changesLogLine)
	specificChangesNumberString := numberRegex.FindString(specificChangesString)
	convertedChanges, err := strconv.Atoi(specificChangesNumberString)

	if err != nil {
		return 0, err
	}

	return convertedChanges, nil
}

/**
 * Parse the git pretty format log line.
 * This is heavily influenced by the format of the pretty format. Look at PrettyFormat for further
 * details on this.
 */
func parsePrettyLogLine(prettyLogLine string) (*CommitData, error) {
	commitData := new(CommitData)
	splitPrettyLogLine := strings.Split(prettyLogLine, PrettyFormatStringSeparator)

	if len(splitPrettyLogLine) != PrettyFormatStringParameterCount() {
		return nil, errors.New("pretty log has an unexpected amount of values")
	}

	authorParsedTime, authorParsedTimeErr := time.Parse(TimeFormat, splitPrettyLogLine[1])
	if authorParsedTimeErr != nil {
		return nil, authorParsedTimeErr
	}

	committerParsedTime, committerParsedTimeErr := time.Parse(TimeFormat, splitPrettyLogLine[4])
	if committerParsedTimeErr != nil {
		return nil, committerParsedTimeErr
	}

	commitData.id = splitPrettyLogLine[0]
	commitData.authorTime = authorParsedTime.Unix()
	commitData.authorName = splitPrettyLogLine[2]
	commitData.authorEmail = splitPrettyLogLine[3]
	commitData.committerTime = committerParsedTime.Unix()
	commitData.committerName = splitPrettyLogLine[5]
	commitData.committerEmail = splitPrettyLogLine[6]

	return commitData, nil
}
