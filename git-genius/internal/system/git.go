package system

import (
	"os"
	"os/exec"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/ui"
)

/* ============================================================
   Core Git Executor
   ============================================================ */

/*
RunGit executes a git command in the selected project directory
and logs errors centrally.
*/
func RunGit(args ...string) error {
	cmd := GitCmd(args...)
	if err := cmd.Run(); err != nil {
		LogError("git "+strings.Join(args, " "), err)
		return err
	}
	return nil
}

/*
GitCmd returns a prepared *exec.Cmd for git,
used when caller needs stdout/stderr control.
*/
func GitCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)

	// Run inside selected WorkDir (if set)
	cfg := config.Load()
	if cfg.WorkDir != "" {
		cmd.Dir = cfg.WorkDir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

/* ============================================================
   Repo Checks
   ============================================================ */

/*
IsGitRepo checks if the selected directory is a git repository
*/
func IsGitRepo() bool {
	cfg := config.Load()

	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if cfg.WorkDir != "" {
		cmd.Dir = cfg.WorkDir
	}

	return cmd.Run() == nil
}

/*
EnsureGitRepo ensures the selected directory is a git repository.
If not, it asks user permission to initialize it.
*/
func EnsureGitRepo() bool {
	if IsGitRepo() {
		return true
	}

	ui.Warn("Selected directory is not a git repository")

	if !ui.Confirm("Do you want to initialize a git repository here?") {
		ui.Error("Git repository required to continue")
		return false
	}

	if err := RunGit("init"); err != nil {
		ui.Error("Failed to initialize git repository")
		return false
	}

	ui.Success("Git repository initialized")
	return true
}