package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

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
