package system

import (
	"os"
	"os/exec"
	"strings"
)

// RunGit executes a git command and logs errors centrally
func RunGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		LogError("git "+strings.Join(args, " "), err)
		return err
	}
	return nil
}
