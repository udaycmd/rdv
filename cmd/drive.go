package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/drives"
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/utils"
)

var (
	add    string
	remove string
	show   bool
)

var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "add or remove a remote drive",
	Long:  "add or remove a remote drive",
	PreRun: func(cmd *cobra.Command, args []string) {
		if add == "" && remove == "" && !show {
			utils.ExitOnError("must provide at least one of the flags for the command")
		} else if add == remove && add != "" {
			utils.ExitOnError("cannot add and remove the same drive")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := internal.LoadCfg()
		if err != nil {
			utils.ExitOnError("%s", err.Error())
		}

		showMe := func() {
			utils.Log(utils.Info, "Supported drives:")
			for _, d := range drives.SupportedDrives {
				fmt.Printf("- %s\n", d.GetCfg().Name)
			}

			for _, d := range cfg.Drives {
				if d.Status == internal.Selected {
					d.GetInfo()
				}
			}
		}

		if add != "" {
			err := driveAdd(add, cfg)
			if err != nil {
				showMe()
				utils.ExitOnError("%s", err.Error())
			}

			utils.Log(utils.Success, "Added %s successfully", add)
		}

		if remove != "" {
			err := driveRemove(remove, cfg)
			if err != nil {
				utils.ExitOnError("%s", err.Error())
			}

			utils.Log(utils.Success, "Removed %s successfully", remove)
		}

		if show {
			showMe()
		}
	},
}

func driveAdd(dn string, c *internal.RdvConfig) error {
	for _, d := range c.Drives {
		if d.Name == dn && d.Status != internal.Revoked {
			return fmt.Errorf("%s already linked with rdv", dn)
		}
	}

	for _, dr := range drives.SupportedDrives {
		dCfg := dr.GetCfg()

		if dCfg.Name == dn {
			c.Drives = append(
				c.Drives,
				internal.DriveProviderConfig{
					Name:   dCfg.Name,
					Id:     dCfg.ClientId,
					Status: internal.Default,
				})
			err := oauth.Authorize(dr)
			if err != nil {
				return err
			}

			return c.SaveCfg()
		}
	}

	return fmt.Errorf("%s is not supported by rdv", dn)
}

func driveRemove(dn string, c *internal.RdvConfig) error {
	for i, d := range c.Drives {
		if d.Name == dn {
			if d.Status != internal.Revoked {
				for _, dr := range drives.SupportedDrives {
					dCfg := dr.GetCfg()

					if dCfg.Name == dn {
						err := oauth.RevokeToken(dr)
						if err != nil {
							return err
						}
						c.Drives[i].Status = internal.Revoked
						return c.SaveCfg()
					}
				}
			} else {
				return fmt.Errorf("%s already disconnected", dn)
			}
		}
	}

	return fmt.Errorf("%s is not connected to rdv", dn)
}

func init() {
	rootCmd.AddCommand(driveCmd)
	driveCmd.Flags().StringVarP(&add, "add", "a", "", "drive to add")
	driveCmd.Flags().StringVarP(&remove, "remove", "r", "", "drive to remove")
	driveCmd.Flags().BoolVarP(&show, "show", "s", false, "show supported drives")
}
