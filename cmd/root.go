package cmd

import (
	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/utils"
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
		utils.ExitOnError("%s", err.Error())
	}
}
