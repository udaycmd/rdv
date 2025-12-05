package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/udaycmd/rdv/utils"
	"golang.org/x/oauth2"
)

var (
	RdvConfFileName string = ".rdv.conf"
	RdvConfFilePath string = filepath.Join(homeDir(), RdvConfFileName)
)

type DriveProviderConfig struct {
	Name string        `json:"name"`
	Id   string        `json:"client_id"`
	T    *oauth2.Token `json:"token"`
}

type RdvCfg struct {
	Ver    string                `json:"ver"`    // cli version
	Drives []DriveProviderConfig `json:"drives"` // drives configuration
}

func homeDir() string {
	res, err := os.UserHomeDir()
	if err != nil {
		utils.ExitOnError("%s", err.Error())
	}
	return res
}

func LoadCfg() (*RdvCfg, error) {
	data, err := os.ReadFile(RdvConfFilePath)
	if os.IsNotExist(err) {
		return &RdvCfg{Ver: utils.Version, Drives: []DriveProviderConfig{}}, nil
	} else if err != nil {
		return nil, err
	}

	cfg := &RdvCfg{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *RdvCfg) SaveCfg() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize configuration data: %w", err)
	}

	return os.WriteFile(RdvConfFilePath, data, 0600)
}
