package logread

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/logformat"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

const splitterString = "__SPLITATME__"

var numberRegex = regexp.MustCompile("[0-9]+")
var insertionsRegex = regexp.MustCompile("([0-9]+) insertions?")
var deletionsRegex = regexp.MustCompile("([0-9]+) deletions?")
var filesChangedRegex = regexp.MustCompile("([0-9]+) files changed?")
var prettyLogLineRegex = regexp.MustCompile(fmt.Sprintf("%s([\\s\\S]*?)%s", logformat.PrettyFormatStringStart, logformat.PrettyFormatStringEnd))

func ParseCommitLog(commitLog string) ([]*common.Commit, error) {
	prettyFormatLines := prettyLogLineRegex.FindAllString(commitLog, -1)
	statLinesString := prettyLogLineRegex.ReplaceAllString(commitLog, splitterString)
	// We want to use the pretty log lines in this case as a separator; if we don't remove the first
	// one then we will have an empty string at the first index of statLines
	statLinesString = strings.Replace(statLinesString, splitterString, "", 1)
	statLines := strings.Split(statLinesString, splitterString)

	prettyFormatLineCount := len(prettyFormatLines)
	statLineCount := len(statLines)

	if prettyFormatLineCount != statLineCount {
		return nil, fmt.Errorf("mismatched pretty log lines (%+v) and stat lines (%+v)", prettyFormatLineCount, statLineCount)
	}

	parsedCommits := make([]*common.Commit, prettyFormatLineCount)

	for i := 0; i < prettyFormatLineCount; i++ {
		commitString := prettyFormatLines[i] + statLines[i]
		parsedCommit, err := ParseCommit(commitString)

		if err != nil {
			return nil, err
		}

		parsedCommits[i] = parsedCommit
	}

	return parsedCommits, nil
}

/**
 * Parse a commit according to the standard format defined in commitformat.go
 *
 * Example commit:

	PRETTYFORMATSTART__4610c5caa1b48f113ee87f48aeace2846a474957__SEPARATOR__Sun, 4 Jun 2023 16:35:34 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Sun, 4 Jun 2023 16:35:34 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Replace use of reflect.DeepEqual with use of new cmp library__SEPARATOR__Signed-off-by: Claudio Cambra <developer@claudiocambra.com>__PRETTYFORMATEND

	go.mod                                                 |   1 +
	go.sum                                                 |   2 ++
	internal/db/testing/sqlite_test.go                     |   8 +++-----
	internal/db/testing/sqlite_test_utils.go               |  12 ++++--------
	pkg/common/changes_test.go                             | 105 ++++++++++++++++++++++++++++++++++++---------------------------------------------------------------------
	pkg/common/email_test.go                               |   8 ++++----
	pkg/logread/commitparse_test.go                        |   8 +++-----
	pkg/statistics/authorgroups/domaingroupsreport_test.go |   6 ++----
	pkg/statistics/authorgroups/domaingroupssql_test.go    |  15 ++++++---------
	9 files changed, 61 insertions(+), 104 deletions(-)

 **/

func ParseCommit(rawCommit string) (*common.Commit, error) {
	commitLogLines := strings.Split(rawCommit, logformat.PrettyFormatStringEnd)
	prettyLogLine := commitLogLines[0]
	changesLogLine := commitLogLines[len(commitLogLines)-1]

	insertions, _ := parseChangesLine(changesLogLine, insertionsRegex)
	deletions, _ := parseChangesLine(changesLogLine, deletionsRegex)
	filesChanged, _ := parseChangesLine(changesLogLine, filesChangedRegex)

	commit, err := parsePrettyLogLine(prettyLogLine)
	if err != nil {
		return nil, err
	}

	commit.NumInsertions = insertions
	commit.NumDeletions = deletions
	commit.NumFilesChanged = filesChanged

	return commit, nil
}

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
func parsePrettyLogLine(prettyLogLine string) (*common.Commit, error) {
	commitData := new(common.Commit)
	splitPrettyLogLine := strings.Split(prettyLogLine, logformat.PrettyFormatStringSeparator)

	expectedParameterCount := logformat.PrettyFormatStringParameterCount()
	prettyLogLineValueCount := len(splitPrettyLogLine)
	if prettyLogLineValueCount != expectedParameterCount {
		return nil, fmt.Errorf("pretty log has an unexpected amount of values: expected %d, received %d", expectedParameterCount, prettyLogLineValueCount)
	}

	// Clean start and end pretty format separators
	firstString := splitPrettyLogLine[0]
	lastString := splitPrettyLogLine[len(splitPrettyLogLine)-1]

	splitPrettyLogLine[0] = strings.Replace(firstString, logformat.PrettyFormatStringStart, "", -1)
	splitPrettyLogLine[len(splitPrettyLogLine)-1] = strings.Replace(lastString, logformat.PrettyFormatStringEnd, "", -1)

	commitData.Id = splitPrettyLogLine[0]

	authorParsedTime, authorParsedTimeErr := time.Parse(common.TimeFormat, splitPrettyLogLine[1])
	if authorParsedTimeErr != nil {
		return nil, authorParsedTimeErr
	}
	commitData.AuthorTime = authorParsedTime.Unix()

	commitData.Author = common.Person{
		Name:  splitPrettyLogLine[2],
		Email: splitPrettyLogLine[3],
	}

	committerParsedTime, committerParsedTimeErr := time.Parse(common.TimeFormat, splitPrettyLogLine[4])
	if committerParsedTimeErr != nil {
		return nil, committerParsedTimeErr
	}
	commitData.CommitterTime = committerParsedTime.Unix()

	commitData.Committer = common.Person{
		Name:  splitPrettyLogLine[5],
		Email: splitPrettyLogLine[6],
	}

	commitData.Subject = splitPrettyLogLine[7]
	commitData.Body = splitPrettyLogLine[8]

	return commitData, nil
}
