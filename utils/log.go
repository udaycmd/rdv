// CLI logging utilities
package utils

import (
	"fmt"
	"os"
)

type Level int

const (
	Info Level = iota
	Success
	Error
	Warn
)

func Log(lvl Level, s string, a ...any) {
	m := fmt.Sprintf(s, a...)

	switch lvl {
	case Info:
		fmt.Printf("%s\n", Colorize(Gray, "[Info] ")+m)
	case Success:
		fmt.Printf("%s\n", Colorize(Green, "[Success] ")+m)
	case Error:
		fmt.Fprintf(os.Stderr, "%s\n", Colorize(Red, "[Error] ")+m)
	case Warn:
		fmt.Printf("%s\n", Colorize(Yellow, "[Warn] ")+m)
	default:
		fmt.Printf("%s\n", m)
	}
}

func ExitOnError(message string, a ...any) {
	Log(Error, message, a...)
	os.Exit(1)
}

func ExitOnSuccess(message string, a ...any) {
	Log(Success, message, a...)
	os.Exit(0)
}
