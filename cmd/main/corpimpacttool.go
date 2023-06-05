package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/logread"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/authorgroups/corpimpact"
)

func main() {
	var (
		ingestDbPath         = flag.String("ingest-db-path", "", "path to database file")
		readDbPath           = flag.String("read-db-path", "", "path to database file")
		repoPath             = flag.String("repo-path", "", "path to git repository")
		domainGroupsFilePath = flag.String("domain-groups-file-path", "", "file containing email domain groups")
	)

	flag.Parse()

	if *ingestDbPath != "" {

		if *repoPath == "" {
			log.Fatalf("Cannot ingest git repository commits to a database file without a path for said file.")
		}

		sqlb := new(db.SQLiteBackend)
		err := sqlb.Open(*ingestDbPath)
		if err != nil {
			log.Fatalf("Error opening sqlite database, received error: %s", err)
			os.Exit(0)
		}

		ingestRepoCommits(*ingestDbPath, *repoPath)
		sqlb.Close()

	} else if *readDbPath != "" && *domainGroupsFilePath != "" {

		if *domainGroupsFilePath == "" {
			log.Println("WARNING: No valid domain groupings file has been provided")
		}

		printDomainGroups(*readDbPath, *domainGroupsFilePath)
	}

	log.Fatalf("No valid individual repo or batch operation specified. Exiting.")
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

	log.Println("Starting commit ingest.")
	sqlb.AddCommits(commits)
	log.Println("Finished ingesting commits!")

	sqlb.Close()
}

func printDomainGroups(readDbPath string, domainGroupsFilePath string) {
	sqlb := new(db.SQLiteBackend)
	err := sqlb.Open(readDbPath)
	if err != nil {
		log.Fatalf("Error opening sqlite database, received error: %s", err)
		os.Exit(0)
	}

	groupsJsonBytes, err := os.ReadFile(domainGroupsFilePath)
	if err != nil {
		log.Fatalf("Error opening domain groups json file: %s", err)
		sqlb.Close()
	}

	var groups map[string][]string
	err = json.Unmarshal(groupsJsonBytes, &groups)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		sqlb.Close()
	}

	corpReport := corpimpact.NewCorporateReport(groups, sqlb, "Corporate")
	corpReport.Generate()
	fmt.Printf("%+v", corpReport)

	sqlb.Close()
}