package internal

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/udaycmd/rdv/utils"
)

type driveStatus int

const (
	Default driveStatus = iota
	Selected
	Revoked
)

var (
	RdvConfFileName string = ".rdv.conf"
	RdvConfFilePath string = filepath.Join(getHomeDir(), RdvConfFileName)
	RdvUserId       string = getUserId()
)

type DriveProviderConfig struct {
	Name   string      `json:"name"`
	Id     string      `json:"client_id"`
	Status driveStatus `json:"status"`
}

type RdvConfig struct {
	Ver    string                `json:"ver"`    // cli version
	Drives []DriveProviderConfig `json:"drives"` // configuration for each drive
}

func getHomeDir() string {
	res, err := os.UserHomeDir()
	if err != nil {
		utils.ExitOnError("%s", err.Error())
	}
	return res
}

func getUserId() string {
	u, err := user.Current()
	if err != nil {
		utils.ExitOnError("%s", err.Error())
	}
	return u.Uid
}

func LoadCfg() (*RdvConfig, error) {
	f, err := os.ReadFile(RdvConfFilePath)
	if os.IsNotExist(err) {
		return &RdvConfig{Ver: utils.Version, Drives: []DriveProviderConfig{}}, nil
	} else if err != nil {
		return nil, err
	}

	cfg := &RdvConfig{}
	if err := json.Unmarshal(f, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (d *DriveProviderConfig) GetInfo() {
	utils.Log(utils.Info, "Provider: %s, Id: %s", d.Name, d.Id)
}

func (c *RdvConfig) SaveCfg() error {
	s, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(RdvConfFilePath, s, 0600)
}
