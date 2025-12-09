package cmd

import (
	"fmt"
	"strings"

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
	use    bool
)

var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "add or remove a remote drive",
	Long:  "add or remove a remote drive",
	PreRun: func(cmd *cobra.Command, args []string) {
		if add == "" && remove == "" && !show && !use {
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

		if add != "" {
			err := driveAdd(add, cfg)
			if err != nil {
				showDrives(cfg)
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

		if use {
			err := useDrive(cfg)
			if err != nil {
				utils.ExitOnError("%s", err.Error())
			}
		}

		if show {
			showDrives(cfg)
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

func useDrive(c *internal.RdvConfig) error {
	switch len(c.Drives) {
	case 0:
		return fmt.Errorf("no connected drives found")
	case 1:
		if c.Drives[0].Status == internal.Revoked {
			utils.Log(utils.Info, "A disconnected drive is found:")
			c.Drives[0].GetInfo()

			utils.Log(utils.Info, "Would you like to connect to it again? (y/n):")
			choice := ""
			fmt.Scanln(&choice)

			if strings.ToLower(choice) == "y" {
				c.Drives[0].Status = internal.Selected

				err := c.SaveCfg()
				if err != nil {
					return err
				}
			}

			utils.Log(utils.Success, "Ok")
		} else {
			c.Drives[0].Status = internal.Selected

			err := c.SaveCfg()
			if err != nil {
				return err
			}

			utils.Log(utils.Success, "Automatically selected the sole drive")
		}

		return nil
	}

	utils.Log(utils.Info, "Select one of the configured drive")
	i := 0
	for {
		if i == len(c.Drives) {
			break
		}

		switch c.Drives[i].Status {
		case internal.Default:
			fmt.Printf("%s (%d)\n", c.Drives[i].Name, i+1)
		case internal.Selected:
			fmt.Printf("%s (%d) *active\n", c.Drives[i].Name, i+1)
		case internal.Revoked:
			fmt.Printf("%s (%d) *disconnected\n", c.Drives[i].Name, i+1)
		}

		i++
	}

	choice := 0
	utils.Log(utils.Info, "Specify your choice (1 - %d)", i)
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return err
	}

	if choice > i || choice < 1 {
		return fmt.Errorf("invalid choice")
	}

	c.Drives[choice-1].Status = internal.Selected
	for i := range c.Drives {
		if i != choice-1 && c.Drives[i].Status == internal.Selected {
			c.Drives[i].Status = internal.Default
		}
	}

	return c.SaveCfg()
}

func showDrives(c *internal.RdvConfig) {
	utils.Log(utils.Info, "Supported drives:")
	for _, d := range drives.SupportedDrives {
		fmt.Printf("- %s\n", d.GetCfg().Name)
	}

	for _, d := range c.Drives {
		if d.Status == internal.Selected {
			d.GetInfo()
		}
	}
}

func init() {
	rootCmd.AddCommand(driveCmd)
	driveCmd.Flags().StringVarP(&add, "add", "a", "", "drive to add")
	driveCmd.Flags().StringVarP(&remove, "remove", "r", "", "drive to remove")
	driveCmd.Flags().BoolVarP(&use, "use", "u", false, "use an added drive")
	driveCmd.Flags().BoolVarP(&show, "show", "s", false, "show supported drives")
}
