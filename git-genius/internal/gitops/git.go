package gitops

import (
	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/* ============================================================
   Helpers
   ============================================================ */

func CurrentBranch() string {
	return config.Load().Branch
}

func CurrentRemote() string {
	return config.Load().Remote
}

/* ============================================================
   Core Git Operations
   ============================================================ */

func Status() {
	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("status"); err != nil {
		ui.Error("Failed to get git status (see error.log)")
	}
}

func Push(msg string) {
	if msg == "" {
		ui.Error("Commit message cannot be empty")
		return
	}

	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("add", "."); err != nil {
		ui.Error("Failed to stage files")
		return
	}

	if err := system.RunGit("commit", "-m", msg); err != nil {
		ui.Warn("Nothing to commit")
		return
	}

	cfg := config.Load()
	if err := system.RunGit("push", cfg.Remote, cfg.Branch); err != nil {
		ui.Error("Push failed (see error.log)")
		return
	}

	ui.Success("Changes pushed successfully")
}

func Pull() {
	if !system.EnsureGitRepo() {
		return
	}

	cfg := config.Load()

	ui.Info("Fetching latest changes...")
	if err := system.RunGit("fetch", cfg.Remote, cfg.Branch); err != nil {
		ui.Error("Fetch failed")
		return
	}

	ui.Info("Merging changes...")
	if err := system.RunGit("merge", cfg.Remote+"/"+cfg.Branch); err != nil {
		ui.Error("Merge conflict detected — resolve manually")
		return
	}

	ui.Success("Pulled latest changes")
}

func Fetch() {
	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("fetch", "--all"); err != nil {
		ui.Error("Fetch failed")
		return
	}
	ui.Success("Fetched all remotes")
}

/* ============================================================
   Branch & Remote
   ============================================================ */

func SwitchBranch() {
	if !system.EnsureGitRepo() {
		return
	}

	name := ui.Input("New branch name")
	if name == "" {
		ui.Error("Branch name cannot be empty")
		return
	}

	if err := system.RunGit("checkout", "-B", name); err != nil {
		ui.Error("Failed to switch branch")
		return
	}

	cfg := config.Load()
	cfg.Branch = name
	config.Save(cfg)

	ui.Success("Switched to branch: " + name)
}

func SwitchRemote() {
	if !system.EnsureGitRepo() {
		return
	}

	name := ui.Input("Remote name")
	url := ui.Input("Remote URL")

	if name == "" || url == "" {
		ui.Error("Remote name and URL are required")
		return
	}

	_ = system.RunGit("remote", "remove", name)

	if err := system.RunGit("remote", "add", name, url); err != nil {
		ui.Error("Failed to add remote")
		return
	}

	cfg := config.Load()
	cfg.Remote = name
	config.Save(cfg)

	ui.Success("Remote set to: " + name)
}
