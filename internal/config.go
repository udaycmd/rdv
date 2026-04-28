package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/udaycmd/rdv/utils"
)

type driveStatus int

const (
	Default driveStatus = iota
	Selected
	Revoked
)

const (
	rdvConfFileName string = ".rdv.config"
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

func LoadCfg() (*RdvConfig, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	f, err := os.ReadFile(filepath.Join(d, rdvConfFileName))
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

func (d *DriveProviderConfig) GetInfo() string {
	return fmt.Sprintf("%s, id: %s", d.Name, d.Id)
}

func (c *RdvConfig) SaveCfg() error {
	d, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	s, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(d, rdvConfFileName), s, 0600)
}

func (c *RdvConfig) GetSelectedDrive() *DriveProviderConfig {
	for i := range c.Drives {
		if c.Drives[i].Status == Selected {
			return &c.Drives[i]
		}
	}

	return nil
}
