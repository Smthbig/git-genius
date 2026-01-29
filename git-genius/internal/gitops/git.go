package gitops

import (
	"os"
	"os/exec"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
run executes git commands and logs errors centrally
*/
func run(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		system.LogError("git "+strings.Join(args, " "), err)
		return err
	}
	return nil
}

/*
Helpers
*/
func CurrentBranch() string {
	return config.Load().Branch
}

func CurrentRemote() string {
	return config.Load().Remote
}

/* ---------- Core Git Operations ---------- */

func Status() {
	if err := run("status"); err != nil {
		ui.Error("Failed to get git status (see error.log)")
	}
}

func Push(msg string) {
	if msg == "" {
		ui.Error("Commit message cannot be empty")
		return
	}

	if err := run("add", "."); err != nil {
		ui.Error("Failed to stage files")
		return
	}

	if err := run("commit", "-m", msg); err != nil {
		ui.Warn("Nothing to commit")
		return
	}

	cfg := config.Load()
	if err := run("push", cfg.Remote, cfg.Branch); err != nil {
		ui.Error("Push failed (see error.log)")
		return
	}

	ui.Success("Changes pushed successfully")
}

func Pull() {
	cfg := config.Load()

	ui.Info("Fetching latest changes...")
	if err := run("fetch", cfg.Remote, cfg.Branch); err != nil {
		ui.Error("Fetch failed")
		return
	}

	ui.Info("Merging changes...")
	if err := run("merge", cfg.Remote+"/"+cfg.Branch); err != nil {
		ui.Error("Merge conflict detected — resolve manually")
		return
	}

	ui.Success("Pulled latest changes")
}

func Fetch() {
	if err := run("fetch", "--all"); err != nil {
		ui.Error("Fetch failed")
		return
	}
	ui.Success("Fetched all remotes")
}

/* ---------- Branch & Remote ---------- */

func SwitchBranch() {
	name := ui.Input("New branch name")
	if name == "" {
		ui.Error("Branch name cannot be empty")
		return
	}

	if err := run("checkout", "-B", name); err != nil {
		ui.Error("Failed to switch branch")
		return
	}

	cfg := config.Load()
	cfg.Branch = name
	config.Save(cfg)

	ui.Success("Switched to branch: " + name)
}

func SwitchRemote() {
	name := ui.Input("Remote name")
	url := ui.Input("Remote URL")

	if name == "" || url == "" {
		ui.Error("Remote name and URL are required")
		return
	}

	run("remote", "remove", name)

	if err := run("remote", "add", name, url); err != nil {
		ui.Error("Failed to add remote")
		return
	}

	cfg := config.Load()
	cfg.Remote = name
	config.Save(cfg)

	ui.Success("Remote set to: " + name)
}

/* ---------- Setup (GUIDED + BEGINNER FRIENDLY) ---------- */

func Setup() {
	cfg := config.Load()

	ui.Header("Git Genius Setup")

	/* Branch */
	if b := ui.Input("Default branch [" + cfg.Branch + "]"); b != "" {
		cfg.Branch = b
	}

	/* Remote */
	if r := ui.Input("Remote [" + cfg.Remote + "]"); r != "" {
		cfg.Remote = r
	}

	/* GitHub Token Help */
	ui.Info("GitHub Token is required for authentication & GitHub API access")
	ui.Info("How to create a token:")
	ui.Info("1. Open: https://github.com/settings/tokens")
	ui.Info("2. Click: Generate new token (classic)")
	ui.Info("3. Note: git-genius")
	ui.Info("4. Select scope: repo")
	ui.Info("5. Generate & COPY the token (shown only once)")

	if !ui.Confirm("Do you want to set / update GitHub token now?") {
		ui.Warn("Skipping GitHub token setup")
		config.Save(cfg)
		return
	}

	token := ui.Input("Paste GitHub token")
	if token == "" {
		ui.Warn("Empty token, skipping")
		config.Save(cfg)
		return
	}

	if err := github.Save(token); err != nil {
		ui.Error("Failed to save token")
		system.LogError("saving github token failed", err)
		return
	}

	user, err := github.Validate()
	if err != nil {
		ui.Error("Invalid GitHub token")
		system.LogError("github token validation failed", err)
		github.Delete()
		ui.Warn("Token removed. Please retry setup.")
		return
	}

	if user == "offline-mode" {
		ui.Warn("Offline mode: token saved but not validated")
	} else {
		ui.Success("GitHub authenticated as: " + user)
	}

	config.Save(cfg)
	ui.Success("Setup completed successfully")
}