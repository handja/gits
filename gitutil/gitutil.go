package gitutil

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func GetGitRepos() []os.FileInfo {
	gitDirectories := make([]os.FileInfo, 0)
	directories, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, directory := range directories {
		if directory.IsDir() && isGitRepo(directory) {
			gitDirectories = append(gitDirectories, directory)
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

func ExecuteGitCommand(directoryName string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = "./" + directoryName
	out, err := cmd.Output()
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		log.Println()
		log.Printf("%s - %s : %s", red("Error"), directoryName, string(out))
		log.Fatal(err)
		return nil, err
	}
	return out, nil
}

func FetchAllBranches(directoryName string) {
	ExecuteGitCommand(directoryName, "fetch", "--all")
}

func GetCurrentBranch(directoryName string) string {
	// git rev-parse --abbrev-ref HEAD
	out := ExecuteGitCommand(directoryName, "rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(string(out))
}

func GetUnpushedBranches(directoryName string) ([]string, []string, []string, bool) {
	localBranches := make([]string, 0)
	remoteBranches := make([]string, 0)
	unpushedBranches := make([]string, 0)
	notUptodateBranches := make([]string, 0)
	aheadBranches := make([]string, 0)
	out := ExecuteGitCommand(directoryName, "branch", "-a")
	allBranches := strings.Split(string(out), "\n")
	allBranches = allBranches[0 : len(allBranches)-1]
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
	out := ExecuteGitCommand(directoryName, "rev-list", branch+"..remotes/origin/"+branch, "--count")
	return strings.TrimSpace(string(out)) != "0"
}

func isAheadBranch(branch string, directoryName string) bool {
	out := ExecuteGitCommand(directoryName, "rev-list", "remotes/origin/"+branch+".."+branch, "--count")
	return strings.TrimSpace(string(out)) != "0"
}

func IsFilesNotStagedOrCommitedOnCurrentBranch(directoryName string) bool {
	out := ExecuteGitCommand(directoryName, "status", "-s")
	return len(strings.TrimSpace(string(out))) > 0
}
