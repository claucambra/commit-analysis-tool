package git

import (
	"fmt"
	"os/exec"
)

func ReadCommits(repoPath string) ([]*CommitData, error) {
	cmd := exec.Command("git",
		"-C", repoPath,
		"log",
		"--no-merges",
		"--branches",
		"--remotes",
		fmt.Sprintf("--pretty=format:%s", PrettyFormatString()),
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
	commits, err := ParseCommitLog(outString)
	if err != nil {
		return nil, err
	}

	return commits, nil
}
