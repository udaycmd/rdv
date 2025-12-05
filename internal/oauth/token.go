// Token serialization and revocation procedures
package oauth

import (
	"fmt"

	"github.com/udaycmd/rdv/internal"
	"golang.org/x/oauth2"
)

func GetToken(p OauthProvider) (*oauth2.Token, error) {
	cfg, err := internal.LoadCfg()
	if err != nil {
		return nil, err
	}

	for _, d := range cfg.Drives {
		if d.Name == p.GetCfg().Name {
			return d.T, nil
		}
	}

	return nil, fmt.Errorf("unknown oauth provider")
}

func SetToken(p OauthProvider, t *oauth2.Token) error {
	cfg, err := internal.LoadCfg()
	if err != nil {
		return err
	}

	for _, d := range cfg.Drives {
		if d.Name == p.GetCfg().Name {
			d.T = t
			return cfg.SaveCfg()
		}
	}

	return fmt.Errorf("unknown oauth provider")
}

// TODO
func Revoke(p OauthProvider) error {
	return nil
}
