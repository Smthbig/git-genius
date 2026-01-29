package menu

import (
	"fmt"
	"os"

	"git-genius/internal/doctor"
	"git-genius/internal/gitops"
	"git-genius/internal/setup"
	"git-genius/internal/ui"
)

func Start() {
	for {
		ui.Clear()
		ui.Header("Git Genius v1.0")

		// Current context
		fmt.Println("Branch :", gitops.CurrentBranch())
		fmt.Println("Remote :", gitops.CurrentRemote())
		fmt.Println()

		// Menu options
		fmt.Println("1) Push changes")
		fmt.Println("2) Pull changes")
		fmt.Println("3) Fetch all remotes")
		fmt.Println("4) Switch branch")
		fmt.Println("5) Switch remote")
		fmt.Println("6) Git status")
		fmt.Println("7) Setup / Reconfigure")
		fmt.Println("8) Doctor (health check)")
		fmt.Println("9) Exit")

		choice := ui.Input("Select option")

		switch choice {
		case "1":
			gitops.Push(ui.Input("Commit message"))

		case "2":
			gitops.Pull()

		case "3":
			gitops.Fetch()

		case "4":
			gitops.SwitchBranch()

		case "5":
			gitops.SwitchRemote()

		case "6":
			gitops.Status()

		case "7":
			setup.Run()

		case "8":
			doctor.Run()

		case "9":
			ui.Info("Goodbye 👋")
			os.Exit(0)

		default:
			ui.Error("Invalid option, please try again")
		}

		ui.Pause()
	}
}
