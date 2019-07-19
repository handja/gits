package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/handja/gits/gitmessage"
	"github.com/handja/gits/gitpull"
	"github.com/handja/gits/gitstatus"
	"github.com/handja/gits/gitutil"
	"github.com/handja/gits/help"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'pull', 'status' or 'help' subcommands")
		os.Exit(1)
	}
	pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
	pullFetch := pullCmd.Bool("fetch", true, "fetch before pull")

	switch os.Args[1] {
	case "pull":
		pullCmd.Parse(os.Args[2:])
		isWithoutFetch := !*pullFetch
		if isWithoutFetch {
			executeAsynchronousGitsCommand(PULL_WITHOUT_FETCH)
		} else if len(os.Args) > 2 {
			fmt.Println("expected '-fetch' subcommands or nothing")
			os.Exit(1)
		} else {
			executeAsynchronousGitsCommand(PULL)
		}
	case "status":
		executeAsynchronousGitsCommand(STATUS)
	case "help":
		help.DisplayHelp()
	case "poule":
		fmt.Println()
		gitmessage.DisplayPocpocMessage()
	default:
		fmt.Println("Run 'gits help' for usage.")
		os.Exit(1)
	}
}

func executeAsynchronousGitsCommand(gitsCommandType GitsCommandType) {
	fmt.Printf("Waiting ...")
	gitDirectories := gitutil.GetGitRepos()
	if len(gitDirectories) == 0 {
		fmt.Println("\rNo git directories")
		os.Exit(1)
	}
	var wg sync.WaitGroup
	wg.Add(len(gitDirectories))
	switch gitsCommandType {
	case STATUS:
		gitstatus.Status(gitDirectories, &wg)
	case PULL:
		gitpull.PullAllBranches(gitDirectories, &wg)
	case PULL_WITHOUT_FETCH:
		gitpull.PullAllBranchesWithoutFetch(gitDirectories, &wg)
	}
	color.Green("Done")
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

type GitsCommandType int

const (
	STATUS = iota
	PULL
	PULL_WITHOUT_FETCH
)

func (g GitsCommandType) String() string {
	return [...]string{"STATUS", "PULL", "PULL_WITHOUT_FETCH"}[g]
}

// git pull on a specified branch
// add verbose option to display all directories
// switch to a specific branch
