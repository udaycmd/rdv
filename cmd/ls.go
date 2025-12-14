package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal/drives"
	// "github.com/udaycmd/rdv/utils"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list contents of a drive",
	Long:  "list contents of a drive",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func lsDrive(drive drives.Drive, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeoutPeriod)
	defer cancel()

	_, err := drive.View(ctx, id)
	if err != nil {
		return err
	}

	// for _, f := range files {
	// 	utils.Log()
	// }

	return nil
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
