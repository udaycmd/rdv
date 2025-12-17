package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/briandowns/spinner"
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
			utils.ExitOnError("%s\n", err.Error())
		}

		Config, err = internal.LoadCfg()
		if err != nil {
			utils.ExitOnError("%s\n", err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		d := Config.GetSelectedDrive()
		if d == nil {
			utils.ExitOnError("No seleted drive found!\n")
		}

		ctx, cancel := context.WithTimeout(context.Background(), RequestTimeoutPeriod)
		defer cancel()

		drive, err := drives.NewDriveFromProvider(ctx, d.Name)
		if err != nil {
			utils.ExitOnError("%s\n", err.Error())
		}

		err = ls(drive, dir)
		if err != nil {
			utils.ExitOnError("%s\n", err.Error())
		}
	},
}

func ls(drive drives.Drive, id string) error {
	s := spinner.New(spinner.CharSets[11], SpinnerDuration)
	s.Color("magenta")
	s.Prefix = utils.Slogf(utils.Info, "Fetching the contents of %s... ", dir)

	s.Start()
	filesMeta, err := drive.View(id)
	s.Stop()
	fmt.Println()

	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 10, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n",
		utils.Colorize(utils.Cyan, "Name"),
		utils.Colorize(utils.Cyan, "Size"),
		utils.Colorize(utils.Cyan, "Id"),
		utils.Colorize(utils.Cyan, "Last Modified"),
	)

	for _, fm := range filesMeta {
		name := fm.Name
		sz := fmt.Sprintf("%d kb", fm.Size/1024)
		id := fm.Id
		modTime := fm.LastModified.Format(time.DateTime)

		print := ""
		if fm.IsDir {
			print = fmt.Sprintf("%s\t%s\t%s\t%s\t\n",
				utils.Colorize(utils.Yellow, name),
				utils.Colorize(utils.Yellow, sz),
				utils.Colorize(utils.Yellow, id),
				utils.Colorize(utils.Yellow, modTime),
			)
		} else {
			print = fmt.Sprintf("%s\t%s\t%s\t%s\t\n",
				utils.Colorize(utils.Green, name),
				utils.Colorize(utils.Green, sz),
				utils.Colorize(utils.Green, id),
				utils.Colorize(utils.Green, modTime),
			)
		}

		fmt.Fprintf(w, "%s", print)
	}

	return w.Flush()
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringVarP(&dir, "dir", "d", "root", "id of the directory")
}
