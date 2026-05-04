package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/utils"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Show version and commit",
	Long:    "Show version and commit",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logger.Logf(log.InfoLevel, "Rdv version: %s, revision: %s", utils.Version, utils.Commit)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
