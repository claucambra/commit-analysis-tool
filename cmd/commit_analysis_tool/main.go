package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/logread"
)

func main() {
	var (
		ingestDbPath = flag.String("ingest-db-path", "", "path to database file")
		readDbPath   = flag.String("read-db-path", "", "path to database file")
		repoPath     = flag.String("repo-path", "", "path to git repository")
		//domainGroupsFile = flag.String("domain-groups-file", "", "file containing email domain groups")
	)

	flag.Parse()

	if *ingestDbPath == "" && *readDbPath == "" {
		log.Fatalf("Must provide a database to ingest to or read from. Quitting.")
		os.Exit(0)
	}

	if *ingestDbPath != "" && *repoPath == "" {
		log.Fatalf("Must provide a git repository path to ingest from. Quitting.")
		os.Exit(0)
	}

	if *ingestDbPath != "" {
		ingestRepoCommits(*ingestDbPath, *repoPath)
	}
}

func ingestRepoCommits(ingestDbPath string, repoPath string) {
	sqlb := new(db.SQLiteBackend)
	err := sqlb.Open(ingestDbPath)

	if err != nil {
		log.Fatalf("Error opening sqlite database, received error: %s", err)
		os.Exit(0)
	}

	err = sqlb.Setup()
	if err != nil {
		log.Fatalf("Error setting up sqlite database, received error: %s", err)
		os.Exit(0)
	}

	commits, err := logread.ReadCommits(repoPath)
	if err != nil {
		log.Fatalf("Error reading commits at %s: %s", repoPath, err)
	}

	sqlb.AddCommits(commits)
	fmt.Printf("Done!\n")
	os.Exit(0)
}
