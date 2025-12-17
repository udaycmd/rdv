package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/utils"
)

var (
	RequestTimeoutPeriod                     = 12 * time.Second
	SpinnerDuration                          = 50 * time.Millisecond
	Config               *internal.RdvConfig = nil
)

var rootCmd = &cobra.Command{
	Use:               "rdv",
	Short:             "Access your cloud drive storage from the terminal!",
	Long:              "rdv (Remote Drive View) is a cli tool that can fetch and upload files and directories to the specified drive.",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	SilenceErrors:     true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.ExitOnError("%s\n", err.Error())
	}
}
