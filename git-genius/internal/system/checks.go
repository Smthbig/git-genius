package system

import (
	"bufio"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"git-genius/internal/ui"
)

var Online bool = false

// -------------------- GIT CHECK & INSTALL --------------------

func EnsureGitInstalled() {
	// ✅ FIX: err declared OUTSIDE if (Go scoping rule)
	_, err := exec.LookPath("git")
	if err == nil {
		return // Git already installed
	}

	// Git not found
	LogError("git not installed", err)
	ui.Error("Git is not installed on this system")

	// Windows: do not auto-install
	if runtime.GOOS == "windows" {
		ui.Warn("Automatic Git install is not supported on Windows")
		ui.Info("Download Git from: https://git-scm.com/downloads")
		os.Exit(1)
	}

	// Ask user for permission
	if !ui.Confirm("Do you want to install Git now?") {
		ui.Error("Git is required to continue")
		os.Exit(1)
	}

	if err := installGit(); err != nil {
		LogError("git install failed", err)
		ui.Error("Automatic Git installation failed")
		ui.Warn("Please install Git manually and retry")
		os.Exit(1)
	}

	ui.Success("Git installed successfully 🎉")
}

// -------------------- INSTALL LOGIC --------------------

func installGit() error {
	switch runtime.GOOS {
	case "linux":
		return installGitLinux()
	case "darwin":
		return installGitMac()
	default:
		return exec.ErrNotFound
	}
}

func installGitLinux() error {
	// Detect package manager
	if exists("apt") {
		return runInstall("sudo apt update && sudo apt install -y git")
	}
	if exists("dnf") {
		return runInstall("sudo dnf install -y git")
	}
	if exists("yum") {
		return runInstall("sudo yum install -y git")
	}
	if exists("pacman") {
		return runInstall("sudo pacman -S --noconfirm git")
	}
	if exists("apk") {
		return runInstall("sudo apk add git")
	}

	return exec.ErrNotFound
}

func installGitMac() error {
	// Try Homebrew
	if exists("brew") {
		return runInstall("brew install git")
	}

	ui.Warn("Homebrew not found")
	ui.Info("Installing Xcode Command Line Tools")
	cmd := exec.Command("xcode-select", "--install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// -------------------- HELPERS --------------------

func exists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func runInstall(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = bufio.NewReader(os.Stdin)
	return cmd.Run()
}

// -------------------- OTHER SYSTEM CHECKS --------------------

func EnsureGitRepo() {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		LogError("not inside git repo", err)
		ui.Error("This is not a git repository")
		os.Exit(1)
	}
}

func CheckInternet() {
	client := http.Client{Timeout: 3 * time.Second}
	_, err := client.Get("https://github.com")
	Online = err == nil
}
