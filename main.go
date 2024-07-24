package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"

	flag "github.com/spf13/pflag"
)

var ReposList string

func main() {
	flagSync := flag.BoolP("sync", "s", false, "Get all the repos in config and sync it using gh repo sync")
	flagForce := flag.BoolP("force", "f", false, "Append the --force flag in gh repo sync")
	flagRepoList := flag.StringP("repolist", "r", filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "syncrepos"), "Provide a repo list to sync")
	flagCreateRepoList := flag.BoolP("createrepolist", "c", false, "Create the repo list if it's not present")

	flag.Parse()

	ReposList = *flagRepoList

	if *flagCreateRepoList {
		createRepoList()
	}

	if flag.NArg() == 0 {
		printRepos()
	} else {
		for _, repo := range flag.Args() {
			addRepo(repo)
		}
	}

	if *flagSync {
		syncRepos(*flagForce)
	}
}

func createRepoList() {
	if _, err := os.Stat(ReposList); os.IsNotExist(err) {
		file, err := os.Create(ReposList)
		if err != nil {
			log.Fatalf("Failed to create repo list file: %v", err)
		}
		defer file.Close()
	}
}

func syncRepos(force bool) {
	repos := readReposFromFile()
	var wg sync.WaitGroup
	errCh := make(chan error, len(repos))

	for repo := range repos {
		wg.Add(1)
		go func(repo string) {
			log.Println("STARTED", repo)
			defer wg.Done()
			defer log.Println("DONE", repo)
			if err := runSyncCommand(repo, force); err != nil {
				errCh <- err
			}
		}(repo)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runSyncCommand(repo string, force bool) error {
	cmd := exec.Command("gh", "repo", "sync", repo)

	if force {
		cmd.Args = append(cmd.Args, "--force")
	}

	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func printRepos() {
	runOnEachLine(func(s string) {
		fmt.Println(s)
	})
}

func addRepo(repo string) {
	if !regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-_]*[a-zA-Z0-9])?\/[a-zA-Z0-9]([a-zA-Z0-9-_]*[a-zA-Z0-9])?$`).MatchString(repo) {
		fmt.Printf("Repositories must be: username/repository but found %s\n", repo)
		return
	}

	repos := readReposFromFile()

	if _, exists := repos[repo]; exists {
		fmt.Println("Repository already exists in the file.")
		return
	}

	writeRepoToFile(repo)
	fmt.Println("Repository added:", repo)
}

func runOnEachLine(run func(string)) {
	readFile, err := os.Open(ReposList)
	if err != nil {
		log.Fatal(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line != "" {
			run(line)
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func readReposFromFile() map[string]struct{} {
	repos := make(map[string]struct{})
	file, err := os.Open(ReposList)
	if err != nil {
		log.Fatalf("Couldn't read file: %v", err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line != "" {
			repos[line] = struct{}{}
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatal(err)
	}

	return repos
}

func writeRepoToFile(repo string) {
	writeFile, err := os.OpenFile(ReposList, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer writeFile.Close()

	if _, err := writeFile.WriteString(repo + "\n"); err != nil {
		log.Fatal(err)
	}
}
