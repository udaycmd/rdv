package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/drives"
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/utils"
	"github.com/zalando/go-keyring"
)

var (
	add    string
	revoke string
	list   bool
	use    bool
)

var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "add or remove a remote drive",
	Long:  "add or remove a remote drive",
	PreRun: func(cmd *cobra.Command, args []string) {
		if add == "" && revoke == "" && !list && !use {
			utils.ExitOnError("must provide at least one of the flags for the command")
		} else if add == revoke && add != "" {
			utils.ExitOnError("cannot add and remove the same drive")
		}

		var err error
		Config, err = internal.LoadCfg()
		if err != nil {
			utils.ExitOnError("%s", err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if add != "" {
			err := driveAdd(add, Config)
			if err != nil {
				utils.ExitOnError("%s", err.Error())
			}

			utils.Log(utils.Success, "Added %s successfully", add)
		}

		if revoke != "" {
			err := driveRevoke(revoke, Config)
			if err != nil {
				utils.ExitOnError("%s", err.Error())
			}

			utils.Log(utils.Success, "Removed %s successfully", revoke)
		}

		if use {
			err := useDrive(Config)
			if err != nil {
				utils.ExitOnError("%s", err.Error())
			}
		}

		if list {
			listDrives(Config)
		}
	},
}

func driveAdd(dn string, c *internal.RdvConfig) error {
	for i, d := range c.Drives {
		if d.Name == dn {
			switch d.Status {
			case internal.Revoked:
				utils.Log(utils.Info, "Reconnecting with %s", d.Name)

				if err := oauth.Authorize(drives.GetDriveOauthProvider(d.Name)); err != nil {
					return err
				}

				c.Drives[i].Status = internal.Default
				return c.SaveCfg()
			default:
				return fmt.Errorf("%s already linked with rdv", dn)
			}
		}
	}

	p := drives.GetDriveOauthProvider(dn)
	if p == nil {
		return fmt.Errorf("%s is not supported by rdv", dn)
	}

	config := p.GetConfig()
	c.Drives = append(c.Drives,
		internal.DriveProviderConfig{
			Name:   config.Name,
			Id:     config.ClientId,
			Status: internal.Default,
		})

	err := oauth.Authorize(p)
	if err != nil {
		return err
	}

	return c.SaveCfg()
}

func driveRevoke(dn string, c *internal.RdvConfig) error {
	for i, d := range c.Drives {
		if d.Name == dn {
			switch d.Status {
			case internal.Revoked:
				return fmt.Errorf("%s already disconnected", dn)
			default:
				p := drives.GetDriveOauthProvider(d.Name)
				err := oauth.RevokeToken(p)
				if err != nil {
					return err
				}

				c.Drives[i].Status = internal.Revoked
				return c.SaveCfg()
			}
		}
	}

	// If token is present in keyring but configuration is empty
	err := oauth.RevokeToken(drives.GetDriveOauthProvider(dn))
	if err == keyring.ErrNotFound {
		return fmt.Errorf("%s is not linked with rdv", dn)
	}

	return err
}

func useDrive(c *internal.RdvConfig) error {
	switch len(c.Drives) {
	case 0:
		return fmt.Errorf("no connected drives found")
	case 1:
		if c.Drives[0].Status == internal.Revoked {
			utils.Log(utils.Info, "A disconnected drive is found: %s", c.Drives[0].GetInfo())

			utils.Log(utils.Info, "Would you like to connect to it again? (y/n):")
			choice := ""
			fmt.Scanln(&choice)

			if strings.ToLower(choice) == "y" {
				if err := driveAdd(c.Drives[0].Name, c); err != nil {
					return err
				}

				c.Drives[0].Status = internal.Selected
				err := c.SaveCfg()
				if err != nil {
					return err
				}
			}

			utils.Log(utils.Success, "ok")
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
			fmt.Printf("%s (%d) [active]\n", c.Drives[i].Name, i+1)
		case internal.Revoked:
			fmt.Printf("%s (%d) [disconnected]\n", c.Drives[i].Name, i+1)
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

	if c.Drives[choice-1].Status == internal.Revoked {
		if err := driveAdd(c.Drives[choice-1].Name, c); err != nil {
			return err
		}
	}

	c.Drives[choice-1].Status = internal.Selected
	for i := range c.Drives {
		if i != choice-1 && c.Drives[i].Status == internal.Selected {
			c.Drives[i].Status = internal.Default
		}
	}

	if err := c.SaveCfg(); err != nil {
		return err
	}

	utils.Log(utils.Success, "Now using %s", c.Drives[choice-1].Name)
	return nil
}

func listDrives(c *internal.RdvConfig) {
	utils.Log(utils.Info, "Supported drives:")
	for _, d := range drives.SupportedDriveProviders {
		fmt.Printf("/> %s\n", d.GetConfig().Name)
	}

	d := c.GetSelectedDrive()
	if d != nil {
		utils.Log(utils.Info, "Active drive: %s", d.GetInfo())
	}
}

func init() {
	driveCmd.Flags().StringVarP(&add, "add", "a", "", "link a remote drive to rdv")
	driveCmd.Flags().StringVarP(&revoke, "revoke", "r", "", "disconnects the drive's client from rdv, but not the drive itself")
	driveCmd.Flags().BoolVarP(&use, "use", "u", false, "use a linked drive")
	driveCmd.Flags().BoolVarP(&list, "list", "l", false, "list supported drives")
	rootCmd.AddCommand(driveCmd)
}
