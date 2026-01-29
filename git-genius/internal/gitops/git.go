package gitops

import (
	"bytes"

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

/*
isWorkingTreeDirty checks if there are uncommitted changes
*/
func isWorkingTreeDirty() bool {
	var out bytes.Buffer
	cmd := system.GitCmd("status", "--porcelain")
	cmd.Stdout = &out
	_ = cmd.Run()
	return out.Len() > 0
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

/* ============================================================
   Stash Operations
   ============================================================ */

func StashSave() {
	if !system.EnsureGitRepo() {
		return
	}

	msg := ui.Input("Stash message (optional)")
	args := []string{"stash", "push"}
	if msg != "" {
		args = append(args, "-m", msg)
	}

	if err := system.RunGit(args...); err != nil {
		ui.Error("Failed to stash changes")
		return
	}

	ui.Success("Changes stashed successfully")
}

func StashList() {
	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("stash", "list"); err != nil {
		ui.Error("Failed to list stashes")
	}
}

func StashPop() {
	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("stash", "pop"); err != nil {
		ui.Error("Failed to apply stash")
		return
	}

	ui.Success("Stash applied successfully")
}

/* ============================================================
   Undo Operations
   ============================================================ */

func UndoLastCommit() {
	if !system.EnsureGitRepo() {
		return
	}

	if !ui.Confirm("Undo last commit? (changes will be kept)") {
		ui.Warn("Undo cancelled")
		return
	}

	if err := system.RunGit("reset", "--soft", "HEAD~1"); err != nil {
		ui.Error("Failed to undo last commit")
		return
	}

	ui.Success("Last commit undone (changes preserved)")
}

/* ============================================================
   Smart Pull (Auto-stash + Pull + Pop)
   ============================================================ */

func SmartPull() {
	if !system.EnsureGitRepo() {
		return
	}

	stashed := false

	if isWorkingTreeDirty() {
		ui.Warn("Uncommitted changes detected")
		if !ui.Confirm("Stash changes and continue pull?") {
			ui.Warn("Smart pull cancelled")
			return
		}

		if err := system.RunGit("stash", "push", "-m", "git-genius-auto-stash"); err != nil {
			ui.Error("Auto-stash failed")
			return
		}
		stashed = true
		ui.Success("Changes auto-stashed")
	}

	cfg := config.Load()
	ui.Info("Pulling latest changes...")
	if err := system.RunGit("pull", cfg.Remote, cfg.Branch); err != nil {
		ui.Error("Pull failed")
		return
	}

	if stashed {
		ui.Info("Restoring stashed changes...")
		if err := system.RunGit("stash", "pop"); err != nil {
			ui.Warn("Failed to auto-apply stash — apply manually")
			return
		}
		ui.Success("Stash restored successfully")
	}

	ui.Success("Smart pull completed")
}