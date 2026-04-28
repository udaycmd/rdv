// Token serialization and revocation procedures
package oauth

import (
	"encoding/json"
	"os/user"

	"github.com/udaycmd/rdv/internal/oauth/providers"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func GetToken(id string) (*oauth2.Token, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	key, err := keyring.Get(id, u.Uid)
	if err != nil {
		return nil, err
	}

	t := &oauth2.Token{}
	if err := json.Unmarshal([]byte(key), t); err != nil {
		return nil, err
	}

	return t, nil
}

func SetToken(id string, t *oauth2.Token) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	key, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return keyring.Set(id, u.Uid, string(key))
}

func RevokeToken(p providers.OauthProvider) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	if p.Revoke() != nil {
		return err
	}

	return keyring.Delete(p.GetConfig().ClientID, u.Uid)
}
