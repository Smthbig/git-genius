package menu

import (
	"fmt"
	"os"

	"git-genius/internal/doctor"
	"git-genius/internal/gitops"
	"git-genius/internal/ui"
)

func Start() {
	for {
		ui.Clear()
		ui.Header("Git Genius v1.0")

		fmt.Println("Branch :", gitops.CurrentBranch())
		fmt.Println("Remote :", gitops.CurrentRemote())
		fmt.Println()

		fmt.Println("1) Push changes")
		fmt.Println("2) Pull changes")
		fmt.Println("3) Fetch all remotes")
		fmt.Println("4) Switch branch")
		fmt.Println("5) Switch remote")
		fmt.Println("6) Git status")
		fmt.Println("7) Setup / Reconfigure")
		fmt.Println("8) Doctor (health check)")
		fmt.Println("9) Exit")

		switch ui.Input("Select option") {
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
			gitops.Setup()
		case "8":
			doctor.Run()
		case "9":
			ui.Info("Goodbye 👋")
			os.Exit(0)
		default:
			ui.Error("Invalid option")
		}

		ui.Pause()
	}
}
