package logread

import (
	"fmt"
	"os/exec"

	"github.com/claucambra/commit-analysis-tool/internal/logformat"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

func ReadCommits(repoPath string, reports []common.Report) ([]*common.CommitData, error) {
	cmd := exec.Command("git",
		"-C", repoPath,
		"log",
		"--no-merges",
		"--branches",
		"--remotes",
		fmt.Sprintf("--pretty=format:%s", logformat.PrettyFormatString()),
		"--reverse",
		"--date-order",
		"HEAD",
		"--stat",
		"--stat-width",
		"999")

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outString := string(out)
	commits, err := ParseCommitLog(outString, reports)
	if err != nil {
		return nil, err
	}

	return commits, nil
}
