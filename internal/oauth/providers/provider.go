package providers

import (
	"golang.org/x/oauth2"
)

// High level contract that should be implemented by every Oauth2 provider
type OauthProvider interface {
	Name() string
	GetConfig() *oauth2.Config
	Revoke() error
}
