package providers

import (
	"github.com/udaycmd/rdv/internal/oauth"
	"golang.org/x/oauth2"
)

type dboxAuthProvider struct{}

var dboxEndpoint = oauth2.Endpoint{
	AuthURL:  "https://www.dropbox.com/oauth2/authorize",
	TokenURL: "https://api.dropboxapi.com/oauth2/token",
}

func NewDboxAuthProvider() *dboxAuthProvider {
	return &dboxAuthProvider{}
}

func (g *dboxAuthProvider) GetCfg() *oauth.BaseConfig {
	return &oauth.BaseConfig{
		Name:     "dbox",
		ClientId: "2qnotffuu8vx1z7",
		Ep:       dboxEndpoint,
	}
}

func (d *dboxAuthProvider) Revoke() error {
	return nil
}
