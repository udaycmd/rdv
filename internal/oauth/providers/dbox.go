package providers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/udaycmd/rdv/internal"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

type dboxAuthProvider struct {
	name string
}

var dboxEndpoint = oauth2.Endpoint{
	AuthURL:  "https://www.dropbox.com/oauth2/authorize",
	TokenURL: "https://api.dropboxapi.com/oauth2/token",
}

func NewDboxAuthProvider() *dboxAuthProvider {
	return &dboxAuthProvider{"dbox"}
}

func (d *dboxAuthProvider) Name() string {
	return d.name
}

func (d *dboxAuthProvider) GetConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: "2qnotffuu8vx1z7",
		Endpoint: dboxEndpoint,
	}
}

func (d *dboxAuthProvider) Revoke() error {
	key, err := keyring.Get(d.GetConfig().ClientID, internal.RdvUserId)
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
