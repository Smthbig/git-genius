package menu

import (
	"fmt"
	"os"
	"path/filepath"

	"git-genius/internal/config"
	"git-genius/internal/doctor"
	"git-genius/internal/gitops"
	"git-genius/internal/setup"
	"git-genius/internal/ui"
)

func Start() {
	for {
		ui.Clear()
		ui.Header("Git Genius v1.0")

		// Load config for context
		cfg := config.Load()

		// Resolve project directory
		projectDir := cfg.WorkDir
		if projectDir == "" {
			projectDir, _ = os.Getwd()
		}

		// ---------------- Context Panel ----------------
		fmt.Println("Project :", filepath.Base(projectDir))
		fmt.Println("Path    :", projectDir)
		fmt.Println("Branch  :", gitops.CurrentBranch())
		fmt.Println("Remote  :", gitops.CurrentRemote())

		if cfg.Owner != "" && cfg.Repo != "" {
			fmt.Println("Repo    :", "https://github.com/"+cfg.Owner+"/"+cfg.Repo)
		}
		fmt.Println()

		// ---------------- Menu ----------------
		fmt.Println("=== Daily Git Operations ===")
		fmt.Println("1) Push changes")
		fmt.Println("2) Pull changes")
		fmt.Println("3) Smart Pull (auto stash)")
		fmt.Println("4) Fetch all remotes")
		fmt.Println("5) Git status")
		fmt.Println()

		fmt.Println("=== Branch / Remote ===")
		fmt.Println("6) Switch branch")
		fmt.Println("7) Switch remote")
		fmt.Println()

		fmt.Println("=== Stash & Undo ===")
		fmt.Println("8) Stash changes")
		fmt.Println("9) Stash list")
		fmt.Println("10) Stash pop")
		fmt.Println("11) Undo last commit")
		fmt.Println()

		fmt.Println("=== Tools ===")
		fmt.Println("12) Setup / Reconfigure")
		fmt.Println("13) Doctor (health check)")
		fmt.Println("14) Exit")

		switch ui.Input("Select option") {

		// ---- Daily ops ----
		case "1":
			gitops.Push(ui.Input("Commit message"))

		case "2":
			gitops.Pull()

		case "3":
			gitops.SmartPull()

		case "4":
			gitops.Fetch()

		case "5":
			gitops.Status()

		// ---- Branch / Remote ----
		case "6":
			gitops.SwitchBranch()

		case "7":
			gitops.SwitchRemote()

		// ---- Stash & Undo ----
		case "8":
			gitops.StashSave()

		case "9":
			gitops.StashList()

		case "10":
			gitops.StashPop()

		case "11":
			gitops.UndoLastCommit()

		// ---- Tools ----
		case "12":
			setup.Run()

		case "13":
			doctor.Run()

		case "14":
			ui.Info("Goodbye 👋")
			os.Exit(0)

		default:
			ui.Error("Invalid option, please try again")
		}

		ui.Pause()
	}
}