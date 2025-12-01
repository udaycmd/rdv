// CLI logging utilities
package utils

import (
	"fmt"
	"os"
)

func Info(message string) {
	fmt.Printf("%s\n", Colorize(Gray, "[Info] ")+message)
}

func Warn(message string) {
	fmt.Printf("%s\n", Colorize(Yellow, "[Warn] ")+message)
}

func Success(message string) {
	fmt.Printf("%s\n", Colorize(Green, "[Success] ")+message)
}

func Error(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", Colorize(Red, "[Error] ")+message)
}

func ExitOnError(message string) {
	Error(message)
	os.Exit(1)
}

func ExitOnSuccess(message string) {
	Success(message)
	os.Exit(0)
}
