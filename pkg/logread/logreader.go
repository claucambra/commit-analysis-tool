package logread

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/claucambra/commit-analysis-tool/internal/logformat"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

func ReadCommits(repoPath string) ([]*common.CommitData, error) {
	log.Println("Running git log.")

	cmd := exec.Command("git",
		"--no-pager",
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
		log.Fatalf("Error running git: %s\n", err)
		return nil, err
	}

	log.Println("Git log printed.")

	outString := string(out)

	log.Println("Starting to parse git log.")
	commits, err := ParseCommitLog(outString)
	if err != nil {
		log.Fatalf("Error during commit log parse: %s\n", err)
		return nil, err
	}

	log.Println("Git log parsing complete.")
	return commits, nil
}
