package providers

import (
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/utils"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type gdriveAuthProvider struct{}

func NewGdriveAuthProvider() *gdriveAuthProvider {
	return &gdriveAuthProvider{}
}

func (g *gdriveAuthProvider) GetCfg() *oauth.BaseConfig {
	return &oauth.BaseConfig{
		Name:     "gdrive",
		ClientId: "593200518603-k0ptna6taq593eiulqnd4vfsk1djh0vl.apps.googleusercontent.com",
		Secret:   "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",
		Scopes:   []string{drive.DriveScope},
		Ep:       google.Endpoint,
	}
}

func (g *gdriveAuthProvider) GetInfo() {
	utils.Log(utils.Info, "Selected Drive client: %s, Client Id: %s", g.GetCfg().Name, g.GetCfg().ClientId)
}

func (g *gdriveAuthProvider) Revoke() error {
	return nil
}
