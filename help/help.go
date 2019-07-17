package help

import "fmt"

func DisplayHelp() {
	fmt.Println("Gits is a tool to manage git multi repositories.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("        gits <command>")
	fmt.Println()
	fmt.Println("The commands are :")
	fmt.Println()
	fmt.Println("        status              get all warnings related to each git repositories")
	fmt.Println("        pull                pull all branches of each git repositories")
	fmt.Println("        pull -fetch=false   pull all branches of each git repositories without fetch (fast)")
	fmt.Println("        poule               display a chicken")
}
