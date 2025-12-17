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

func Slogf(lvl Level, t string, a ...any) string {
	m := fmt.Sprintf(t, a...)
	var s string

	switch lvl {
	case Info:
		s = Colorize(Cyan, "Info: ") + m
	case Success:
		s = Colorize(Green, "Success: ") + m
	case Error:
		s = Colorize(Red, "Error: ") + m
	case Warn:
		s = Colorize(Yellow, "Warn: ") + m
	default:
		s = m
	}

	return s
}

func ExitOnError(message string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s", Slogf(Error, message, a...))
	os.Exit(1)
}

func ExitOnSuccess(message string, a ...any) {
	fmt.Fprintf(os.Stdout, "%s", Slogf(Success, message, a...))
	os.Exit(0)
}
