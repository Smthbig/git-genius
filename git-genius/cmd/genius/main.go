package main

import (
	"git-genius/internal/menu"
	"git-genius/internal/system"
)

func main() {
	system.EnsureGitInstalled()
	system.EnsureGitRepo()
	system.CheckInternet()
	menu.Start()
}
