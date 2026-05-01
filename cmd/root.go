package cmd

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/utils"
)

var (
	requestTimeoutPeriod                     = 10 * time.Second
	config               *internal.RdvConfig = nil
)

var rootCmd = &cobra.Command{
	Use:               "rdv",
	Short:             "Access your cloud drive storage from the terminal!",
	Long:              "rdv (Remote Drive View) is a cli tool that can fetch and upload files and directories to the specified drive.",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	SilenceErrors:     true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Logger.Logf(log.FatalLevel, "%s", err.Error())
	}
}

func init() {
	utils.InitLogger(false)
	var err error
	config, err = internal.LoadCfg()
	if err != nil {
		utils.Logger.Fatal("can't load config file", "message", err.Error())
	}
}
