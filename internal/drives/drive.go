package drives

import (
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/internal/oauth/providers"
)

var SupportedDriveProviders = []oauth.OauthProvider{
	providers.NewGdriveAuthProvider(),
	providers.NewDboxAuthProvider(),
}

func GetDriveOauthProvider(name string) oauth.OauthProvider {
	for _, p := range SupportedDriveProviders {
		if p.GetConfig().Name == name {
			return p
		}
	}

	return nil
}

type Drive interface {
	// View the contents of a directory,
	// if id is empty string then root directory of the drive is selected
	View(id string) error

	// Get/download an object with an id
	// from the chosen drive
	Get(id string) error

	// Set/upload an object with an id
	// with an optional name in the chosen drive
	Set(id string, name string) error

	// Deletes an object with an id from the chosen drive
	Delete(id string) error
}
