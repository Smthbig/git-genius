package doctor

import (
	"os"
	"os/exec"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

// Run performs full system + git health check
func Run() {
	ui.Header("Git Genius Doctor 🩺")

	checkGit()
	checkRepo()
	checkGitConfig()
	checkInternet()
	checkGitHubToken()
	checkErrorLog()

	ui.Success("Doctor check completed")
}

// ---------------- CHECKS ----------------

func checkGit() {
	if _, err := exec.LookPath("git"); err != nil {
		ui.Error("Git not installed")
	} else {
		ui.Success("Git installed")
	}
}

func checkRepo() {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		ui.Error("Not inside a git repository")
	} else {
		ui.Success("Inside a git repository")
	}
}

func checkGitConfig() {
	name := gitConfig("user.name")
	email := gitConfig("user.email")

	if name == "" {
		ui.Warn("git user.name not set")
	} else {
		ui.Success("git user.name: " + name)
	}

	if email == "" {
		ui.Warn("git user.email not set")
	} else {
		ui.Success("git user.email: " + email)
	}

	cfg := config.Load()
	ui.Info("Default branch: " + cfg.Branch)
	ui.Info("Default remote: " + cfg.Remote)
}

func checkInternet() {
	if system.Online {
		ui.Success("Internet connection available")
	} else {
		ui.Warn("Offline mode detected")
	}
}

func checkGitHubToken() {
	token := github.Get()
	if token == "" {
		ui.Warn("GitHub token not configured")
		return
	}

	user, err := github.Validate()
	if err != nil {
		ui.Error("GitHub token invalid")
		return
	}

	if user == "offline-mode" {
		ui.Warn("GitHub token validation skipped (offline)")
		return
	}

	ui.Success("GitHub authenticated as: " + user)
}

func checkErrorLog() {
	if _, err := os.Stat(".git/.genius/error.log"); err == nil {
		ui.Warn("Error log exists (.git/.genius/error.log)")
	} else {
		ui.Success("No error log found")
	}
}

// ---------------- HELPERS ----------------

func gitConfig(key string) string {
	out, err := exec.Command("git", "config", "--get", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
