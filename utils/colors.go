// CLI text manipulation utilities
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Gray   = "\033[90m"
)

func Colorize(color, text string) string {
	return color + text + Reset
}

func ClearScreen() error {
	var err error
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin", "freebsd":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		err = fmt.Errorf("unsupported OS")
	}

	if err != nil {
		return err
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
