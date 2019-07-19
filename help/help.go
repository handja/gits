package help

import "fmt"

func DisplayHelp() {
	fmt.Println("Gits is a tool to manage multi git repositories.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("        gits <command>")
	fmt.Println()
	fmt.Println("The commands are :")
	fmt.Println()
	fmt.Println("        status              get all warnings related to each git repositorie")
	fmt.Println("        pull                pull all branches of each git repositorie")
	fmt.Println("        pull -fetch=false   pull all branches of each git repositorie without fetch (faster)")
	fmt.Println("        poule               display a chicken")
}
