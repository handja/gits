package gitpull

import (
	"os"
	"sync"

	"github.com/handja/gits/gitutil"
)

func PullAllBranchesWithoutFetch() {

}

func PullAllBranches(gitDirectories []os.FileInfo, wg *sync.WaitGroup) {
	for _, gitDirectory := range gitDirectories {
		go pullOnRepository(gitDirectory, wg)
	}
	wg.Wait()
}

func pullOnRepository(gitDirectory os.FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	isCurrentBranchWorkingInProgress := gitutil.IsFilesNotStagedOrCommitedOnCurrentBranch(gitDirectory.Name())
	directoryName := gitDirectory.Name()
	if !isCurrentBranchWorkingInProgress {
		gitutil.FetchAllBranches(directoryName)
		_, notUptodateBranches, _, _ := gitutil.GetUnpushedBranches(directoryName)
		currentBranch := gitutil.GetCurrentBranch(directoryName)
		if len(notUptodateBranches) > 0 {
			for _, notUptodateBranch := range notUptodateBranches {
				pullOnBranch(directoryName, notUptodateBranch)
			}
			gitutil.ExecuteGitCommand(directoryName, "checkout", currentBranch)
		}
	}
}

func pullOnBranch(directoryName string, branchName string) {
	checkoutBranch(directoryName, branchName)
	pullGitRepository(directoryName)
}

func checkoutBranch(directoryName string, branchName string) {
	gitutil.ExecuteGitCommand(directoryName, "checkout", branchName)
}

func pullGitRepository(directoryName string) {
	gitutil.ExecuteGitCommand(directoryName, "pull", "--rebase=true")
}
