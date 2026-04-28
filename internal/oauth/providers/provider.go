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

var oauthRegistry = make(map[string]OauthProvider)

func Register(p OauthProvider) {
	oauthRegistry[p.Name()] = p
}

func Get(name string) OauthProvider {
	return oauthRegistry[name]
}

func GetAll() []OauthProvider {
	var providers []OauthProvider
	for _, p := range oauthRegistry {
		providers = append(providers, p)
	}

	return providers
}
