package drives

import (
	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/internal/oauth/providers"
)

var SupportedDrives = []oauth.OauthProvider{
	providers.NewGdriveAuthProvider(),
	providers.NewDboxAuthProvider(),
}

type Drive interface {
	// Get/download an object with an id
	// from the chosen drive
	Get(id string) error

	// Set/upload an object with an id
	// with an optional name in the chosen drive
	Set(id string, name string) error

	// Deletes an object with an id from the chosen drive
	Delete(id string) error
}
