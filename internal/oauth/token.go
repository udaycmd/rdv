// Token serialization and revocation procedures
package oauth

import (
	"encoding/json"

	"github.com/udaycmd/rdv/internal"
	"github.com/udaycmd/rdv/internal/oauth/providers"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func GetToken(id string) (*oauth2.Token, error) {
	key, err := keyring.Get(id, internal.RdvUserId)
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
	key, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return keyring.Set(id, internal.RdvUserId, string(key))
}

func RevokeToken(p providers.OauthProvider) error {
	// server side cleanup
	err := p.Revoke()
	if err != nil {
		return err
	}

	return keyring.Delete(p.GetConfig().ClientID, internal.RdvUserId)
}
