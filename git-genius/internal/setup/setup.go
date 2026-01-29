package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
Run executes the full guided setup
*/
func Run() {
	ui.Header("Git Genius Setup")

	// STEP 0: Select project directory
	if !selectWorkDir() {
		return
	}

	// STEP 1: Ensure git repo (ask + init if missing)
	if !system.EnsureGitRepo() {
		return
	}

	cfg := config.Load()

	// STEP 2: Basic git config
	setupGitBasics(&cfg)

	// STEP 3: GitHub repo selection
	if !setupRepo(&cfg) {
		return
	}

	// STEP 4: GitHub auth + remote
	if !setupGitHubToken(&cfg) {
		return
	}

	config.Save(cfg)

	ui.Header("Setup Summary")
	ui.Success("Project Dir : " + cfg.WorkDir)
	ui.Success("Repository  : https://github.com/" + cfg.Owner + "/" + cfg.Repo)
	ui.Success("Remote      : " + cfg.Remote)
	ui.Success("Branch      : " + cfg.Branch)
	ui.Success("Setup completed successfully 🎉")
}

/* ============================================================
   STEP 0: Project directory selection
   ============================================================ */

func selectWorkDir() bool {
	cfg := config.Load()

	cwd, _ := os.Getwd()
	ui.Info("Current directory: " + cwd)

	if !ui.Confirm("Do you want to use a DIFFERENT project directory?") {
		cfg.WorkDir = cwd
		config.Save(cfg)
		return true
	}

	dir := ui.Input("Enter full path of project directory")
	if dir == "" {
		ui.Error("Directory path cannot be empty")
		return false
	}

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		ui.Error("Invalid directory path")
		return false
	}

	cfg.WorkDir = dir
	config.Save(cfg)

	ui.Success("Project directory set to: " + dir)
	return true
}

/* ============================================================
   STEP 1: Git basics
   ============================================================ */

func setupGitBasics(cfg *config.Config) {
	if b := ui.Input("Default branch [" + cfg.Branch + "]"); b != "" {
		cfg.Branch = b
	}

	if r := ui.Input("Remote name [" + cfg.Remote + "]"); r != "" {
		cfg.Remote = r
	}
}

/* ============================================================
   STEP 2: GitHub repo selection
   ============================================================ */

func setupRepo(cfg *config.Config) bool {
	ui.Header("GitHub Repository")

	// Suggest repo name from folder
	if cfg.Repo == "" {
		base := filepath.Base(cfg.WorkDir)
		if base != "" {
			cfg.Repo = base
		}
	}

	if cfg.Owner == "" {
		cfg.Owner = ui.Input("GitHub username or organisation")
	}

	if cfg.Repo == "" {
		cfg.Repo = ui.Input("Repository name")
	}

	if cfg.Owner == "" || cfg.Repo == "" {
		ui.Error("Owner and repository name are required")
		return false
	}

	ui.Info(fmt.Sprintf(
		"Target repo: https://github.com/%s/%s",
		cfg.Owner, cfg.Repo,
	))

	return true
}

/* ============================================================
   STEP 3: GitHub token + remote
   ============================================================ */

func setupGitHubToken(cfg *config.Config) bool {
	ui.Header("GitHub Authentication")

	ui.Info("GitHub Token is required for HTTPS authentication")
	ui.Info("How to create a token:")
	ui.Info("1. Open: https://github.com/settings/tokens")
	ui.Info("2. Generate new token (classic)")
	ui.Info("3. Note: git-genius")
	ui.Info("4. Select scope: repo")
	ui.Info("5. Copy token")

	if !ui.Confirm("Do you want to configure GitHub token now?") {
		ui.Warn("Skipping GitHub token setup")
		return true
	}

	token := ui.Input("Paste GitHub token")
	if token == "" {
		ui.Warn("Empty token, skipping")
		return true
	}

	if err := github.Save(token); err != nil {
		system.LogError("saving token failed", err)
		ui.Error("Failed to save token")
		return false
	}

	user, err := github.Validate()
	if err != nil {
		system.LogError("token validation failed", err)
		ui.Error("Invalid GitHub token")
		github.Delete()
		return false
	}

	if user != "offline-mode" {
		ui.Success("GitHub authenticated as: " + user)
	}

	// Warn before overwriting remote
	if remoteExists(cfg.Remote) {
		if !ui.Confirm("Remote already exists. Overwrite it?") {
			ui.Warn("Keeping existing remote")
			return true
		}
	}

	if err := configureRemoteWithToken(cfg, token); err != nil {
		system.LogError("remote config failed", err)
		ui.Error("Failed to configure git remote")
		return false
	}

	ui.Success("Git remote configured with token")
	return true
}

/* ============================================================
   Helpers
   ============================================================ */

func configureRemoteWithToken(cfg *config.Config, token string) error {
	url := fmt.Sprintf(
		"https://%s@github.com/%s/%s.git",
		token,
		cfg.Owner,
		cfg.Repo,
	)

	_ = system.RunGit("remote", "remove", cfg.Remote)
	return system.RunGit("remote", "add", cfg.Remote, url)
}

func remoteExists(name string) bool {
	return system.RunGit("remote", "get-url", name) == nil
}
