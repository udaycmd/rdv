package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/log"
)

func OpenFile(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return OpenURL("file://" + path)
}

// [https://github.com/pkg/browser/blob/master/browser.go]
func OpenURL(url string) error {
	var provider string
	providers := []string{"xdg-open", "open", "x-www-browser", "www-browser"}

	switch runtime.GOOS {
	case "linux", "darwin", "freebsd":
		for _, p := range providers {
			if _, err := exec.LookPath(p); err == nil {
				provider = p
				break
			}
		}
	case "windows":
		provider = ""
	default:
		return fmt.Errorf("browser: unsupported OS")
	}

	if provider == "" {
		Logger.Logf(log.InfoLevel, "Please open the following url in your system web browser: \n\n %s", url)
		return nil
	}

	cmd := exec.Command(provider, url)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
