package main

import (
	"github.com/udaycmd/rdv/cmd"
	"github.com/udaycmd/rdv/utils"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		utils.Logger.Fatal("cmd error", "reason", err.Error())
	}
}
