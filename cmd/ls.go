package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/drives"
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/utils"
)

var dir string

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list contents of a drive",
	Long:  "list contents of a drive",
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if err = utils.ClearScreen(); err != nil {
			utils.ExitOnError("%s", err.Error())
		}

		Config, err = internal.LoadCfg()
		if err != nil {
			utils.ExitOnError("%s", err.Error())
		}

		d := Config.GetSelectedDrive()
		if d == nil {
			utils.ExitOnError("No seleted drive found!")
		}

		t, err := oauth.GetToken(d.Id)
		if err != nil {
			utils.ExitOnError("%s", err.Error())
		}

		_ = t
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := ls(nil, dir)
		if err != nil {
			utils.ExitOnError("%s", err.Error())
		}
	},
}

func ls(drive drives.Drive, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeoutPeriod)
	defer cancel()

	_, err := drive.View(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringVarP(&dir, "dir", "d", "", "display the content of the provided")
}
