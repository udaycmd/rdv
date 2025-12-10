package providers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/zalando/go-keyring"
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

func (g *dboxAuthProvider) GetConfig() *oauth.BaseConfig {
	return &oauth.BaseConfig{
		Name:     "dbox",
		ClientId: "2qnotffuu8vx1z7",
		Ep:       dboxEndpoint,
	}
}

func (d *dboxAuthProvider) Revoke() error {
	key, err := keyring.Get(d.GetConfig().ClientId, internal.RdvUserId)
	if err != nil {
		return err
	}

	t := &oauth2.Token{}
	if err := json.Unmarshal([]byte(key), t); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/auth/token/revoke", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+t.AccessToken)
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to revoke token, %s", res.Status)
	}

	return res.Body.Close()
}
