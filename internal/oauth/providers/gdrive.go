package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/user"
	"strings"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type gdriveAuthProvider struct {
	name string
}

func NewGdriveAuthProvider() *gdriveAuthProvider {
	return &gdriveAuthProvider{"gdrive"}
}

func (g *gdriveAuthProvider) Name() string {
	return g.name
}

func (g *gdriveAuthProvider) GetConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "51115061200-h4istsfp5uunp6r0qveu0nk3lbp2pl4h.apps.googleusercontent.com",
		ClientSecret: "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",
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveScope},
	}
}

func (g *gdriveAuthProvider) Revoke() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	key, err := keyring.Get(g.GetConfig().ClientID, u.Uid)
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

func init() {
	Register(NewGdriveAuthProvider())
}
