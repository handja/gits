package gitpull

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/handja/gits/gitutil"
)

func Pull() {
	fmt.Printf("Waiting ...")
	gitDirectories := gitutil.GetGitRepos()
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

func PullAllBranches() {
	fmt.Printf("Waiting ...")
	gitDirectories := gitutil.GetGitRepos()
	if len(gitDirectories) == 0 {
		fmt.Println("\rNo git directories")
		return
	}
	for _, gitDirectory := range gitDirectories {
		isCurrentBranchWorkingInProgress := gitutil.IsFilesNotStagedOrCommitedOnCurrentBranch(gitDirectory.Name())
		directoryName := gitDirectory.Name()
		if !isCurrentBranchWorkingInProgress {
			gitutil.FetchAllBranches(directoryName)
			_, notUptodateBranches, _, _ := gitutil.GetUnpushedBranches(directoryName)
			currentBranch := gitutil.GetCurrentBranch(directoryName)
			if len(notUptodateBranches) > 0 {
				for _, notUptodateBranch := range notUptodateBranches {
					gitutil.ExecuteGitCommand(directoryName, "checkout", notUptodateBranch)
					pullGitRepository(gitDirectory.Name())
				}
				gitutil.ExecuteGitCommand(directoryName, "checkout", currentBranch)
			}
		}
	}
	fmt.Printf("\r\n")
	color.Green("Done")
}

func pullGitRepository(directoryName string) {
	gitutil.ExecuteGitCommand(directoryName, "pull", "--rebase=true")
}
