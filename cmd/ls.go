package cmd

import "github.com/spf13/cobra"

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list contents of a drive",
	Long:  "list contents of a drive",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
