// Token serialization and revocation procedures
package oauth

import (
	"encoding/json"

	"github.com/udaycmd/rdv/internal"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func GetToken(p OauthProvider) (*oauth2.Token, error) {
	key, err := keyring.Get(p.GetCfg().ClientId, internal.RdvUserId)
	if err != nil {
		return nil, err
	}

	t := &oauth2.Token{}
	if err := json.Unmarshal([]byte(key), t); err != nil {
		return nil, err
	}

	return t, nil
}

func SetToken(p OauthProvider, t *oauth2.Token) error {
	key, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return keyring.Set(p.GetCfg().ClientId, internal.RdvUserId, string(key))
}

func RevokeToken(p OauthProvider) error {
	// server side cleanup
	err := p.Revoke()
	if err != nil {
		return err
	}

	return keyring.Delete(p.GetCfg().ClientId, internal.RdvUserId)
}
