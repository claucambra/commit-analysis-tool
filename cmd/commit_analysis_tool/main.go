package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/logread"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/authorgroups"
)

func main() {
	var (
		ingestDbPath         = flag.String("ingest-db-path", "", "path to database file")
		readDbPath           = flag.String("read-db-path", "", "path to database file")
		repoPath             = flag.String("repo-path", "", "path to git repository")
		domainGroupsFilePath = flag.String("domain-groups-file-path", "", "file containing email domain groups")
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
	} else if *readDbPath != "" {
		printDomainGroups(*readDbPath, *domainGroupsFilePath)
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
	sqlb.Close()
	os.Exit(0)
}

func printDomainGroups(readDbPath string, domainGroupsFilePath string) {
	if domainGroupsFilePath == "" {
		log.Fatalf("Cannot create author domain group report without domain groups. Quitting.")
		os.Exit(0)
	}

	sqlb := new(db.SQLiteBackend)
	err := sqlb.Open(readDbPath)
	if err != nil {
		log.Fatalf("Error opening sqlite database, received error: %s", err)
		os.Exit(0)
	}

	groupsJsonBytes, err := ioutil.ReadFile(domainGroupsFilePath)
	if err != nil {
		log.Fatalf("Error opening domain groups json file: %s", err)
		os.Exit(0)
	}

	var groups map[string][]string
	err = json.Unmarshal(groupsJsonBytes, &groups)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		os.Exit(0)
	}

	report := authorgroups.NewDomainGroupsReport(groups)
	report.Generate(sqlb)

	for groupName, _ := range groups {
		fmt.Printf("%+v", report.GroupData(groupName))
	}

	sqlb.Close()
	os.Exit(0)
}
