package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/utils"
)

var (
	requestTimeoutPeriod                     = 10 * time.Second
	config               *internal.RdvConfig = nil
)

var RootCmd = &cobra.Command{
	Use:               "rdv",
	Short:             "Access your cloud drive storage from the terminal!",
	Long:              "rdv (Remote Drive View) is a cli tool that can fetch and upload files and directories to the specified drive.",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	SilenceErrors:     true,
}

func init() {
	utils.InitLogger(false)
	var err error
	config, err = internal.LoadCfg()
	if err != nil {
		utils.Logger.Fatal("can't load config file", "reason", err.Error())
	}
}
