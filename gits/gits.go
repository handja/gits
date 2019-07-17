package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/handja/gits/gitmessage"
	"github.com/handja/gits/gitpull"
	"github.com/handja/gits/gitstatus"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'pull' or 'status' subcommands")
		os.Exit(1)
	}
	pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
	pullAll := pullCmd.Bool("all", false, "pull on all branches")

	switch os.Args[1] {
	case "pull":
		pullCmd.Parse(os.Args[2:])
		if *pullAll {
			gitpull.PullAllBranches()
		} else if len(os.Args) > 2 {
			fmt.Println("expected '-all' subcommands or nothing")
			os.Exit(1)
		} else {
			gitpull.Pull()
		}
	case "status":
		gitstatus.Status()
	case "poule":
		fmt.Println()
		gitmessage.DisplayPocpocMessage()
	default:
		fmt.Println("expected 'pull' or 'status' subcommands")
		os.Exit(1)
	}
}

// add argument for the path
// display the current branch commit and stage
// git pull on each branch (master and develop)
// git pull on a specified branch
// add verbose option to display all directories
// sort result

// "Already up to date." message for git pull if nothing happens
