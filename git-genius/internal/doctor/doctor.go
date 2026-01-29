package doctor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

// Run performs full system + git health check
func Run() {
	ui.Header("Git Genius Doctor 🩺")

	checkGitInstalled()
	checkWorkDir()
	checkGitRepo()
	checkGitConfig()
	checkInternet()
	checkGitHubToken()
	checkErrorLog()

	ui.Success("Doctor check completed")
}

/* ============================================================
   CHECKS
   ============================================================ */

func checkGitInstalled() {
	if _, err := exec.LookPath("git"); err != nil {
		ui.Error("Git not installed")
	} else {
		ui.Success("Git installed")
	}
}

func checkWorkDir() {
	cfg := config.Load()

	if cfg.WorkDir == "" {
		cwd, _ := os.Getwd()
		ui.Info("Project directory (current): " + cwd)
		return
	}

	info, err := os.Stat(cfg.WorkDir)
	if err != nil || !info.IsDir() {
		ui.Error("Configured project directory not found: " + cfg.WorkDir)
		return
	}

	ui.Success("Project directory: " + cfg.WorkDir)
}

func checkGitRepo() {
	if system.IsGitRepo() {
		ui.Success("Git repository detected")
	} else {
		ui.Warn("No git repository found in project directory")
		ui.Info("Run Setup to initialize git repository")
	}
}

func checkGitConfig() {
	name := gitConfig("user.name")
	email := gitConfig("user.email")

	if name == "" {
		ui.Warn("git user.name not set")
		ui.Info("Run: git config --global user.name \"Your Name\"")
	} else {
		ui.Success("git user.name: " + name)
	}

	if email == "" {
		ui.Warn("git user.email not set")
		ui.Info("Run: git config --global user.email \"you@example.com\"")
	} else {
		ui.Success("git user.email: " + email)
	}

	cfg := config.Load()
	ui.Info("Default branch : " + cfg.Branch)
	ui.Info("Default remote : " + cfg.Remote)

	if cfg.Owner != "" && cfg.Repo != "" {
		ui.Info("GitHub repo     : https://github.com/" + cfg.Owner + "/" + cfg.Repo)
	}
}

func checkInternet() {
	if system.Online {
		ui.Success("Internet connection available")
	} else {
		ui.Warn("Offline mode detected")
		ui.Info("GitHub validation & push may fail")
	}
}

func checkGitHubToken() {
	token := github.Get()
	if token == "" {
		ui.Warn("GitHub token not configured")
		ui.Info("Run Setup to configure GitHub authentication")
		return
	}

	user, err := github.Validate()
	if err != nil {
		ui.Error("GitHub token invalid or expired")
		ui.Info("Run Setup to reconfigure token")
		return
	}

	if user == "offline-mode" {
		ui.Warn("GitHub token validation skipped (offline)")
		return
	}

	ui.Success("GitHub authenticated as: " + user)
}

func checkErrorLog() {
	cfg := config.Load()

	base := cfg.WorkDir
	if base == "" {
		base, _ = os.Getwd()
	}

	logPath := filepath.Join(base, ".git", ".genius", "error.log")

	if _, err := os.Stat(logPath); err == nil {
		ui.Warn("Error log exists")
		ui.Info("Check: " + logPath)
	} else {
		ui.Success("No error log found")
	}
}

/* ============================================================
   HELPERS
   ============================================================ */

func gitConfig(key string) string {
	cmd := exec.Command("git", "config", "--get", key)

	cfg := config.Load()
	if cfg.WorkDir != "" {
		cmd.Dir = cfg.WorkDir
	}

	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}