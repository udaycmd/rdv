package main

import (
	"github.com/charmbracelet/log"
	"github.com/udaycmd/rdv/cmd"
	"github.com/udaycmd/rdv/utils"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		utils.Logger.Logf(log.FatalLevel, "%s", err.Error())
	}
}
