package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/logread"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/authorgroups/corpimpact"
)

func main() {
	var (
		batchRead            = flag.String("batch-read", "", "path to file of git clone urls to analyse")
		clonePath            = flag.String("clone-path", "", "path to store cloned repositories in")
		ingestDbPath         = flag.String("ingest-db-path", "", "path to database file")
		readDbPath           = flag.String("read-db-path", "", "path to database file")
		repoPath             = flag.String("repo-path", "", "path to git repository")
		domainGroupsFilePath = flag.String("domain-groups-file-path", "", "file containing email domain groups")
	)

	flag.Parse()

	if *batchRead != "" {

		if *clonePath == "" {
			log.Fatalf("Received empty clone path, don't know where to store cloned repos")
		} else if *domainGroupsFilePath == "" {
			log.Println("WARNING: No valid domain groupings file has been provided")
		}

		batchCloneAndRead(*batchRead, *clonePath, *domainGroupsFilePath)

	} else if *ingestDbPath != "" {

		if *repoPath == "" {
			log.Fatalf("Cannot ingest git repository commits to a database file without a path for said file.")
		}

		sqlb := newSql(*ingestDbPath)
		ingestRepoCommits(*ingestDbPath, *repoPath, sqlb)
		sqlb.Close()

	} else if *readDbPath != "" && *domainGroupsFilePath != "" {

		if *domainGroupsFilePath == "" {
			log.Println("WARNING: No valid domain groupings file has been provided")
		}

		sqlb := newSql(*readDbPath)
		report := generateCorpReport(*readDbPath, *domainGroupsFilePath, sqlb)
		sqlb.Close()

		fmt.Printf("%+v", report)

	} else {

		log.Fatalf("No valid individual repo or batch operation specified. Exiting.")

	}
}

func newSql(dbpath string) *db.SQLiteBackend {
	sqlb := new(db.SQLiteBackend)
	err := sqlb.Open(dbpath)
	if err != nil {
		log.Fatalf("Error opening sqlite database, received error: %s", err)
	}

	return sqlb
}

func ingestRepoCommits(ingestDbPath string, repoPath string, sqlb *db.SQLiteBackend) {
	err := sqlb.Setup()
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
}

func generateCorpReport(readDbPath string, domainGroupsFilePath string, sqlb *db.SQLiteBackend) *corpimpact.CorporateReport {
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

	return corpReport
}

func cloneRepos(urls []string, clonePath string) ([]string, []string) {
	fullClonedPaths := make([]string, len(urls))
	repoNames := make([]string, len(urls))

	for i, url := range urls {
		log.Printf("About to clone git repository: %s", url)

		splitUrl := strings.Split(url, "/")
		repoNameDotGit := splitUrl[len(splitUrl)-1]
		repoName := strings.TrimSuffix(repoNameDotGit, ".git")
		repoNames[i] = repoName

		fullClonePath := filepath.Join(clonePath, repoName)
		fullClonedPaths[i] = fullClonePath

		var cmd *exec.Cmd
		if _, err := os.Stat(fullClonePath); os.IsNotExist(err) {
			cmd = exec.Command("git",
				"clone",
				"--progress",
				url,
				fullClonePath)
		} else {
			cmd = exec.Command("git",
				"-C",
				fullClonePath,
				"pull")
		}

		var stdBuffer bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &stdBuffer)

		cmd.Stdout = mw
		cmd.Stderr = mw

		log.Printf("About to run: %s", cmd.String())

		if err := cmd.Run(); err != nil {
			log.Panic(err)
		}

		log.Printf("Clone of %s now complete.", repoName)
	}

	return fullClonedPaths, repoNames
}

func batchCloneAndRead(urlsJsonFile string, clonePath string, domainGroupsFilePath string) {
	urlsJsonBytes, err := os.ReadFile(urlsJsonFile)
	if err != nil {
		log.Fatalf("Error opening batch fetch urls JSON file: %s", err)
	}

	var urls []string
	err = json.Unmarshal(urlsJsonBytes, &urls)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	fullCsvPath := filepath.Join(clonePath, "corpreport.csv")
	fullCsvFile, err := os.Create(fullCsvPath)
	if err != nil {
		log.Fatalf("Could not create report file: %s", err)
	}

	csvWriter := csv.NewWriter(fullCsvFile)
	firstLineWritten := false

	clonedPaths, repoNames := cloneRepos(urls, clonePath)

	for i, clonedRepoPath := range clonedPaths {
		repoName := repoNames[i]

		log.Printf("\n\nBeginning analysis of: %s\n\n", repoName)

		ingestDbPath := filepath.Join(clonePath, repoName+".db")

		sqlb := newSql(ingestDbPath)

		log.Printf("Beginning commit ingest at %s", ingestDbPath)
		ingestRepoCommits(ingestDbPath, clonedRepoPath, sqlb)
		log.Printf("Commit ingest for %s now complete.", repoName)

		log.Printf("Beginning corporate impact analysis.")
		report := generateCorpReport(ingestDbPath, domainGroupsFilePath, sqlb)

		sqlb.Close()

		fmt.Printf("\n%+v\n", report)

		csvline := report.CSVString(repoName, !firstLineWritten)
		firstLineWritten = true
		err = csvWriter.WriteAll(csvline)
		if err != nil {
			log.Fatalf("Error writing to csv: %s", err)
		}

		// Do CSV file for changes
		repoChangesCsvPath := filepath.Join(clonePath, repoName+"-changes.csv")
		repoChangesCsvFile, err := os.Create(repoChangesCsvPath)
		if err != nil {
			log.Fatalf("Could not create repo changes csv file: %s", err)
		}

		repoChangesDataCSV := report.CSVChangesString(repoName)
		repoChangesWriter := csv.NewWriter(repoChangesCsvFile)
		err = repoChangesWriter.WriteAll(repoChangesDataCSV)
		if err != nil {
			log.Fatalf("Error writing to changes csv: %s", err)
		}

		// Do CSV file for survival
		repoSurvivalCsvPath := filepath.Join(clonePath, repoName+"-survival.csv")
		repoSurvivalCsvFile, err := os.Create(repoSurvivalCsvPath)
		if err != nil {
			log.Fatalf("Could not create repo survival csv file: %s", err)
		}

		repoSurvivalDataCSV := report.CSVSurvivalString(repoName)
		repoSurvivalWriter := csv.NewWriter(repoSurvivalCsvFile)
		err = repoSurvivalWriter.WriteAll(repoSurvivalDataCSV)
		if err != nil {
			log.Fatalf("Error writing to survival csv: %s", err)
		}
	}
}
