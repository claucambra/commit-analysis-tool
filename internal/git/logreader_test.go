package git

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

const repoUrl = "https://github.com/claucambra/commit-analysis-tool.git"
const testDirName = "logreader_test_dir"

func setupTestRepo() (string, error) {
	testDir, err := os.MkdirTemp("", testDirName)
	if err != nil {
		return "", err
	}

	fmt.Printf("Setting up test git environment at %s\n", testDir)
	cmd := exec.Command("git",
		"clone",
		repoUrl,
		testDir)

	cmdErr := cmd.Run()
	return testDir, cmdErr
}

func cleanupTestRepo(path string) {
	os.RemoveAll(path)
}

func TestReadCommits(t *testing.T) {
	repoPath, err := setupTestRepo()
	cleanup := func() { cleanupTestRepo(repoPath) }
	t.Cleanup(cleanup)

	if err != nil {
		t.Fatalf("Test setup failed: %s", err)
	}

	commits, err := ReadCommits(repoPath)
	if err != nil {
		t.Fatalf("Error reading commits: %s", err)
	}

	numReadCommits := len(commits)
	if numReadCommits == 0 {
		t.Fatalf("Read no commits %d", numReadCommits)
	}

	t.Logf("Read %d commits", numReadCommits)
}
