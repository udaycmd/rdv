package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type gdriveAuthProvider struct{}

func NewGdriveAuthProvider() *gdriveAuthProvider {
	return &gdriveAuthProvider{}
}

func (g *gdriveAuthProvider) GetConfig() *oauth.BaseConfig {
	return &oauth.BaseConfig{
		Name:     "gdrive",
		ClientId: "593200518603-k0ptna6taq593eiulqnd4vfsk1djh0vl.apps.googleusercontent.com",
		Secret:   "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",
		Scopes:   []string{drive.DriveScope},
		Ep:       google.Endpoint,
	}
}

func (g *gdriveAuthProvider) Revoke() error {
	key, err := keyring.Get(g.GetConfig().ClientId, internal.RdvUserId)
	if err != nil {
		return err
	}

	t := &oauth2.Token{}
	if err := json.Unmarshal([]byte(key), t); err != nil {
		return err
	}

	body := url.Values{}
	body.Set("token", t.AccessToken)

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/revoke", strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
