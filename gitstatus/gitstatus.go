package gitstatus

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/handja/gits/gitutil"
)

type Repository struct {
	Name                 string
	CurrentBranch        string
	HasWarning           bool
	IsWorkingInProgress  bool
	IsNotOnDevelopBranch bool
	AheadBranches        []string
	UnpushedBranches     []string
	NotUptodateBranches  []string
}

func Status() {
	fmt.Printf("Waiting ...")
	gitDirectories := gitutil.GetGitRepos()
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
		if repository.HasWarning {
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
	gitutil.FetchAllBranches(directoryName)
	unpushedBranches, notUptodateBranches, aheadBranches, isGitFlow := gitutil.GetUnpushedBranches(directoryName)
	currentBranch := gitutil.GetCurrentBranch(directoryName)
	isOnDevelopBranch := (strings.Contains(currentBranch, "develop") && isGitFlow) || (strings.Contains(currentBranch, "master") && !isGitFlow)
	isCurrentBranchWorkingInProgress := gitutil.IsFilesNotStagedOrCommitedOnCurrentBranch(directoryName)
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
