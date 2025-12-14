package utils

import (
	"fmt"
	"time"
)

func Spinner[T any](worker func() (T, error), message string) (T, error) {
	stop := make(chan struct{})
	stopped := make(chan struct{})

	go func() {
		defer close(stopped)

		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		// hide and show cursor on exit
		fmt.Print("\033[?25l")
		defer fmt.Print("\033[?25h")
		i := 0
		for {
			select {
			case <-stop:
				fmt.Print("\033[2K\r")
				return
			case <-ticker.C:
				fmt.Printf("\r%s %s", frames[i%len(frames)], message)
				i++
			}
		}
	}()

	res, err := worker()
	close(stop)

	// prevents immediate return
	<-stopped

	return res, err
}
