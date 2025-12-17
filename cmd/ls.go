package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/drives"
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		d := Config.GetSelectedDrive()
		if d == nil {
			utils.ExitOnError("No seleted drive found!")
		}

		ctx, cancel := context.WithTimeout(context.Background(), RequestTimeoutPeriod)
		defer cancel()

		drive, err := drives.NewDriveFromProvider(ctx, d.Name)
		if err != nil {
			utils.ExitOnError("%s", err.Error())
		}

		err = ls(drive, dir)
		if err != nil {
			utils.ExitOnError("%s", err.Error())
		}
	},
}

func ls(drive drives.Drive, id string) error {
	filesMeta, err := drive.View(id)
	if err != nil {
		return err
	}

	for _, fm := range filesMeta {
		utils.Log(utils.Success, "id: %s", fm.Name)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringVarP(&dir, "dir", "d", "", "display the content of the provided")
}
