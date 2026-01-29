package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
ANSI color codes
*/
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[1;31m"
	Green   = "\033[1;32m"
	Yellow  = "\033[1;33m"
	Blue    = "\033[1;34m"
	Cyan    = "\033[1;36m"
	Magenta = "\033[1;35m"
)

/*
Input helpers
*/
func Input(label string) string {
	fmt.Print(Cyan + label + ": " + Reset)
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return strings.TrimSpace(sc.Text())
}

func SecretInput(label string) string {
	fmt.Print(Cyan + label + ": " + Reset)
	byteInput, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')
	return strings.TrimSpace(string(byteInput))
}

func Confirm(question string) bool {
	for {
		fmt.Print(Yellow + question + " (y/n): " + Reset)
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		ans := strings.ToLower(strings.TrimSpace(sc.Text()))

		if ans == "y" || ans == "yes" {
			return true
		}
		if ans == "n" || ans == "no" {
			return false
		}
		fmt.Println(Red + "Please enter y or n." + Reset)
	}
}

/*
Screen helpers
*/
func Pause() {
	fmt.Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func Header(title string) {
	fmt.Println(Magenta + "========================================" + Reset)
	fmt.Println(Bold + Cyan + " " + title + Reset)
	fmt.Println(Magenta + "========================================" + Reset)
}

/*
Message helpers
*/
func Info(msg string) {
	fmt.Println(Cyan + "ℹ " + msg + Reset)
}

func Success(msg string) {
	fmt.Println(Green + "✔ " + msg + Reset)
}

func Warn(msg string) {
	fmt.Println(Yellow + "⚠ " + msg + Reset)
}

func Error(msg string) {
	fmt.Println(Red + "✘ " + msg + Reset)
}
