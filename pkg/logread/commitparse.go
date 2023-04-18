package logread

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/logformat"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

var numberRegex = regexp.MustCompile("[0-9]+")
var commitRegex = regexp.MustCompile(fmt.Sprintf("^[0-9a-f]+%s", logformat.PrettyFormatStringSeparator))
var insertionsRegex = regexp.MustCompile("([0-9]+) insertions?")
var deletionsRegex = regexp.MustCompile("([0-9]+) deletions?")
var filesChangedRegex = regexp.MustCompile("([0-9]+) files changed?")
var emptyNewLineRegex = regexp.MustCompile(`\n\s*\n`)

func ParseCommitLog(commitLog string, reports []common.Report) ([]*common.CommitData, error) {
	splitCommitLog := emptyNewLineRegex.Split(commitLog, -1)
	numCommits := len(splitCommitLog)
	parsedCommits := make([]*common.CommitData, numCommits)

	for i := 0; i < numCommits; i++ {
		commitString := splitCommitLog[i]
		parsedCommit, err := ParseCommit(commitString)

		if err != nil {
			return nil, err
		}

		parsedCommits[i] = parsedCommit

		for _, report := range reports {
			report.AddCommit(*parsedCommit)
		}
	}

	return parsedCommits, nil
}

/**
 * Parse a commit according to the standard format defined in commitformat.go
 *
 * Example commit:

	1c915e7dd147d4b060c2c241bb966d6f6c6ecde9__SEPARATOR__Sat, 8 Apr 2023 17:47:43 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Wed, 12 Apr 2023 23:21:43 +0000__SEPARATOR__Jean-Baptiste Kempf__SEPARATOR__jb@videolan.org
	modules/gui/macosx/library/VLCLibraryWindow.h                            |  6 +++---
	modules/gui/macosx/library/VLCLibraryWindowPersistentPreferences.h       | 22 +++++++++-------------
	modules/gui/macosx/library/VLCLibraryWindowPersistentPreferences.m       | 30 +++++++++++++++---------------
	modules/gui/macosx/library/audio-library/VLCLibraryAudioViewController.m |  4 ++--
	modules/gui/macosx/library/media-source/VLCMediaSourceBaseDataSource.m   |  4 ++--
	modules/gui/macosx/library/video-library/VLCLibraryVideoViewController.m |  2 +-
	6 files changed, 32 insertions(+), 36 deletions(-)

 **/

func ParseCommit(rawCommit string) (*common.CommitData, error) {
	commitLogLines := strings.Split(rawCommit, "\n")
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
func parsePrettyLogLine(prettyLogLine string) (*common.CommitData, error) {
	commitData := new(common.CommitData)
	splitPrettyLogLine := strings.Split(prettyLogLine, logformat.PrettyFormatStringSeparator)

	if len(splitPrettyLogLine) != logformat.PrettyFormatStringParameterCount() {
		return nil, errors.New("pretty log has an unexpected amount of values")
	}

	authorParsedTime, authorParsedTimeErr := time.Parse(common.TimeFormat, splitPrettyLogLine[1])
	if authorParsedTimeErr != nil {
		return nil, authorParsedTimeErr
	}

	committerParsedTime, committerParsedTimeErr := time.Parse(common.TimeFormat, splitPrettyLogLine[4])
	if committerParsedTimeErr != nil {
		return nil, committerParsedTimeErr
	}

	commitData.Id = splitPrettyLogLine[0]
	commitData.AuthorTime = authorParsedTime.Unix()
	commitData.AuthorName = splitPrettyLogLine[2]
	commitData.AuthorEmail = splitPrettyLogLine[3]
	commitData.CommitterTime = committerParsedTime.Unix()
	commitData.CommitterName = splitPrettyLogLine[5]
	commitData.CommitterEmail = splitPrettyLogLine[6]

	return commitData, nil
}
