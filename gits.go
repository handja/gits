package main

import (
	"fmt"
	"log"
	"os/exec"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"github.com/fatih/color"
	"sort"
)

type Repository struct {
	Name                    string
	CurrentBranch           string
	HasWarning              bool
	IsWorkingInProgress     bool
	IsNotOnDevelopBranch    bool
	AheadBranches           []string
	UnpushedBranches        []string
	NotUptodateBranches     []string
}

func main() {
	displayTitle()
	fmt.Println()
	displayPocpocMessage()
	fmt.Println()

	if len(os.Args) < 2 {
        fmt.Println("expected 'pull' or 'status' subcommands")
        os.Exit(1)
	}

	switch os.Args[1] {
		case "pull":
			pull()
		case "status":
			status()
		default:
			fmt.Println("expected 'pull' or 'status' subcommands")
			os.Exit(1)
	}
}

func displayTitle() {
	fmt.Println("--------------------------")
	fmt.Println("| GITS - multi repo tool |")
	fmt.Println("--------------------------")
}

func displayPocpocMessage() {
	color.Red(" MM")
	fmt.Println("<O \\___/|")
	fmt.Println("  \\_  _/")
	color.Yellow("    ][  O")
}

func pull() {
	fmt.Printf("Waiting ...")
	gitDirectories := getGitRepos()
	if len(gitDirectories) == 0 {
		fmt.Println("\rNo git directories")
		return
	}
	for _, gitDirectory := range gitDirectories {
		pullGitRepository(gitDirectory.Name())
	}
	fmt.Printf("\r\n")
	color.Green("Done")
}

func pullGitRepository(directoryName string) {
	executeGitCommand(directoryName, "pull", "--rebase=true")
}

func status() {
	fmt.Printf("Waiting ...")
	gitDirectories := getGitRepos()
	if len(gitDirectories) == 0 {
		fmt.Println("\rNo git directories")
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(gitDirectories))
	var repositories []Repository
	for _, gitDirectory := range gitDirectories {
		go addGitRepositoryData(gitDirectory.Name(), &repositories, &wg)
	}
	wg.Wait()
	sort.SliceStable(repositories, func(i, j int) bool { return repositories[i].Name < repositories[j].Name })
	displayGitRepositoryWarning(repositories)
	fmt.Println()
	color.Green("Done")
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func displayGitRepositoryWarning(repositories []Repository) {
	isNoWarnings := true
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, repository := range repositories {
		if (repository.HasWarning) {
			if isNoWarnings {
				fmt.Printf("\rDirectories list :")
				fmt.Println()
				fmt.Println()
			}
			fmt.Printf("%s %s\n----------\n", yellow("O"), repository.Name)
			isNoWarnings = false
			if repository.IsNotOnDevelopBranch {
				fmt.Println("current branch (when is not develop branch) : " + repository.CurrentBranch)
			}
			if repository.IsWorkingInProgress {
				color.Red("changes not staged or committed")
			}
			if len(repository.AheadBranches) > 0 {
				color.Blue("ahead branches (commits not pushed) : ")
				fmt.Println("- " + strings.TrimSuffix(strings.Join(repository.AheadBranches, "\n- "), "\n- "))
			}
			if len(repository.UnpushedBranches) > 0 {
				fmt.Println("unpushed branches : ")
				fmt.Println("- " + strings.TrimSuffix(strings.Join(repository.UnpushedBranches, "\n- "), "\n- "))
			}
			if len(repository.NotUptodateBranches) > 0 {
				fmt.Println("not up-to-date branches : ")
				fmt.Println("- " + strings.TrimSuffix(strings.Join(repository.NotUptodateBranches, "\n- "), "\n- "))
			}
			fmt.Println()
		}
	}

	if isNoWarnings {
		fmt.Printf("\rNothing to report.")
	}
}

func addGitRepositoryData(directoryName string, repositories *[]Repository, wg *sync.WaitGroup) {
	defer wg.Done()
	var repository = Repository{Name: directoryName}
	fetchAllBranches(directoryName)
	unpushedBranches, notUptodateBranches, aheadBranches, isGitFlow := getUnpushedBranches(directoryName)
	currentBranch := getCurrentBranch(directoryName)
	isOnDevelopBranch := (strings.Contains(currentBranch, "develop") && isGitFlow) || (strings.Contains(currentBranch, "master") && !isGitFlow)
	isCurrentBranchWorkingInProgress := isFilesNotStagedOrCommitedOnCurrentBranch(directoryName)
	if len(unpushedBranches) > 0 || !isOnDevelopBranch || len(notUptodateBranches) > 0 || len(aheadBranches) > 0 || isCurrentBranchWorkingInProgress {
		repository.HasWarning = true
		repository.CurrentBranch = currentBranch
		repository.IsNotOnDevelopBranch = !isOnDevelopBranch
		repository.IsWorkingInProgress = isCurrentBranchWorkingInProgress 
		repository.AheadBranches = aheadBranches
		repository.UnpushedBranches = unpushedBranches
		repository.NotUptodateBranches = notUptodateBranches
	}
	*repositories = append(*repositories, repository)
}

func getGitRepos() []os.FileInfo {
	gitDirectories := make([]os.FileInfo, 0)
	directories, err := ioutil.ReadDir(".")
	if err != nil {
        log.Fatal(err)
	}
	for _, directory := range directories {
		if directory.IsDir() && isGitRepo(directory) {
			gitDirectories = append(gitDirectories, directory);
		}
	} 
	return gitDirectories
}

func isGitRepo(directory os.FileInfo) bool {
	isGitRepo := false
	directories, err := ioutil.ReadDir("./" + directory.Name())
	if err != nil {
        log.Fatal(err)
	}
	for _, directory := range directories {
		if directory.IsDir() && directory.Name() == ".git" {
			isGitRepo = true
		}
	} 
	return isGitRepo
}

func fetchAllBranches(directoryName string) {
	executeGitCommand(directoryName, "fetch", "--all")
}

func getCurrentBranch(directoryName string) string {
	// git rev-parse --abbrev-ref HEAD
	out := executeGitCommand(directoryName, "rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(string(out))
}

func getUnpushedBranches(directoryName string) ([]string, []string, []string, bool) {
	localBranches := make([]string, 0)
	remoteBranches := make([]string, 0)
	unpushedBranches := make([]string, 0)
	notUptodateBranches := make([]string, 0)
	aheadBranches := make([]string, 0)
	out := executeGitCommand(directoryName, "branch", "-a")
	allBranches := strings.Split(string(out), "\n")
	allBranches = allBranches[0:len(allBranches) - 1]
	for _, branch := range allBranches {
		if strings.Contains(branch, "remotes/origin") {
			remoteBranch := strings.Replace(branch, "remotes/origin/", "", 1)
			remoteBranch = strings.TrimSpace(remoteBranch)
			remoteBranches = append(remoteBranches, remoteBranch)
		} else {
			localBranch := strings.Replace(branch, "*", "", 1)
			localBranch = strings.TrimSpace(localBranch)
			localBranches = append(localBranches, localBranch)
		}
	}
	for _, localBranch := range localBranches {
		if isUnpushedBranch(localBranch, remoteBranches) {
			unpushedBranches = append(unpushedBranches, localBranch)
		} else {
			if isNotUptodateBranch(localBranch, directoryName) {
				notUptodateBranches = append(notUptodateBranches, localBranch)
			}
			if isAheadBranch(localBranch, directoryName) {
				aheadBranches = append(aheadBranches, localBranch)
			}
		}
	}
	isGitFlow := containsDevelopBranch(remoteBranches)
	return unpushedBranches, notUptodateBranches, aheadBranches, isGitFlow
}

func isUnpushedBranch(branch string, remoteBranches []string) bool {
	for _, remoteBranch := range remoteBranches {
		if remoteBranch == branch {
			return false
		}
	}
	return true
}

func containsDevelopBranch(remoteBranches []string) bool {
	for _, remoteBranch := range remoteBranches {
		if remoteBranch == "develop" {
			return true
		}
	}
	return false
}

func isNotUptodateBranch(branch string, directoryName string) bool {
	// git rev-list develop..remotes/origin/develop --count
	out := executeGitCommand(directoryName, "rev-list", branch + "..remotes/origin/" + branch, "--count")
	return strings.TrimSpace(string(out)) != "0"
}

func isAheadBranch(branch string, directoryName string) bool {
	out := executeGitCommand(directoryName, "rev-list", "remotes/origin/" + branch + ".." + branch, "--count")
	return strings.TrimSpace(string(out)) != "0"
}

func isFilesNotStagedOrCommitedOnCurrentBranch(directoryName string) bool {
	out := executeGitCommand(directoryName, "status", "-s")
	return len(strings.TrimSpace(string(out))) > 0
}

func executeGitCommand(directoryName string, args ...string) []byte {
	cmd := exec.Command("git", args...)
	cmd.Dir = "./" + directoryName
	out, err := cmd.Output()
	if err != nil {
		log.Print(err.Error())
		log.Fatal(err)
	}
	return out
}

// add argument for the path
// display the current branch commit and stage
// git pull on each branch (master and develop)
// git pull on a specified branch
// add verbose option to display all directories
// sort result

// "Already up to date." message for git pull if nothing happens
